package profiler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gitlab.com/dotpe/mindbenders/profiler/entity"
	"gitlab.com/dotpe/mindbenders/profiler/uploader"

	pprofile "github.com/google/pprof/profile"
)

// outChannelSize specifies the size of the profile output channel.
const outChannelSize = 5

var (
	mu             sync.Mutex
	activeProfiler *profiler
)

// Start starts the profiler. It may return an error if an API key is not provided by means of
// the WithAPIKey option, or if a hostname is not found.
func Start(opts ...Option) error {
	mu.Lock()
	defer mu.Unlock()
	if activeProfiler != nil {
		activeProfiler.stop()
	}
	p, err := newProfiler(opts...)
	if err != nil {
		return err
	}
	activeProfiler = p
	activeProfiler.run()
	return nil
}

// Stop cancels any ongoing profiling or upload operations and returns after
// everything has been stopped.
func Stop() {
	mu.Lock()
	if activeProfiler != nil {
		activeProfiler.stop()
		activeProfiler = nil
	}
	mu.Unlock()
}

// profiler collects and sends preset profiles to the dotpe API at a given frequency
// using a given configuration.
type profiler struct {
	cfg      *Config                           // profile configuration
	out      chan entity.Batch                 // upload queue
	exit     chan struct{}                     // exit signals the profiler to stop; it is closed after stopping
	stopOnce sync.Once                         // stopOnce ensures the profiler is stopped exactly once.
	wg       sync.WaitGroup                    // wg waits for all goroutines to exit when stopping.
	met      *metrics                          // metric collector state
	prev     map[ProfileType]*pprofile.Profile // previous collection results for delta profiling
	uploader uploader.IProfileUploader
}

// newProfiler creates a new, unstarted profiler.
func newProfiler(opts ...Option) (*profiler, error) {
	cfg, err := defaultConfig()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// uploadTimeout defaults to DefaultUploadTimeout, but in theory a user might
	// set it to 0 or a negative value. However, it's not clear what this should
	// mean, and most meanings we could assign seem to be bad: Not having a
	// timeout is dangerous, having a timeout that fires immediately breaks
	// uploading, and silently defaulting to the default timeout is confusing.
	// So let's just stay clear of all of this by not allowing such values.
	//
	// see similar discussion: https://github.com/golang/go/issues/39177
	if cfg.UploadTimeout <= 0 {
		return nil, fmt.Errorf("invalid upload timeout, must be > 0: %s", cfg.UploadTimeout)
	}
	for pt := range cfg.Types {
		if _, ok := profileTypes[pt]; !ok {
			return nil, fmt.Errorf("unknown profile type: %d", pt)
		}
	}

	p := profiler{
		cfg:  cfg,
		out:  make(chan entity.Batch, outChannelSize),
		exit: make(chan struct{}),
		met:  newMetrics(),
		prev: make(map[ProfileType]*pprofile.Profile),
	}
	p.uploader = cfg.uploader
	return &p, nil
}

// run runs the profiler.
func (p *profiler) run() {

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		tick := time.NewTicker(p.cfg.Period)
		defer tick.Stop()
		p.met.reset(now()) // collect baseline metrics at profiler start
		p.collect(tick.C)
	}()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.send()
	}()
}

// collect runs the profile types found in the configuration whenever the ticker receives
// an item.
func (p *profiler) collect(ticker <-chan time.Time) {
	defer close(p.out)
	for {
		select {
		case <-ticker:
			now := now()
			bat := entity.Batch{
				Start: now,
				// NB: while this is technically wrong in that it does not
				// record the actual start and end timestamps for the batch,
				// it is how the backend understands the client-side
				// configured CPU profile duration: (start-end).
				End: now.Add(p.cfg.CpuDuration),
			}

			for _, t := range p.enabledProfileTypes() {
				profs, err := p.runProfile(t)
				if err != nil {
					continue
				}
				for _, prof := range profs {
					bat.AddProfile(prof)
				}
			}
			p.enqueueUpload(bat)
		case <-p.exit:
			return
		}
	}
}

// enabledProfileTypes returns the enabled profile types in a deterministic
// order. The CPU profile always comes first because people might spot
// interesting events in there and then try to look for the counter-part event
// in the mutex/heap/block profile. Deterministic ordering is also important
// for delta profiles, otherwise they'd cover varying profiling periods.
func (p *profiler) enabledProfileTypes() []ProfileType {
	order := []ProfileType{
		CPUProfile,
		HeapProfile,
		GoroutineProfile,
		MetricsProfile,
	}
	enabled := []ProfileType{}
	for _, t := range order {
		if _, ok := p.cfg.Types[t]; ok {
			enabled = append(enabled, t)
		}
	}
	return enabled
}

// enqueueUpload pushes a batch of profiles onto the queue to be uploaded. If there is no room, it will
// evict the oldest profile to make some. Typically a batch would be one of each enabled profile.
func (p *profiler) enqueueUpload(bat entity.Batch) {
	for {
		select {
		case p.out <- bat:
			return // 👍
		default:
			// queue is full; evict oldest
			select {
			case <-p.out:
			default:
				// this case should be almost impossible to trigger, it would require a
				// full p.out to completely drain within nanoseconds or extreme
				// scheduling decisions by the runtime.
			}
		}
	}
}

// send takes profiles from the output queue and uploads them.
func (p *profiler) send() {
	for {
		select {
		case <-p.exit:
			return
		case bat := <-p.out:
			p.UploadProfile(&bat)
		}
	}
}

func (p *profiler) UploadProfile(b *entity.Batch) error {
	for _, prof := range b.Profiles {
		// 	// 	fileDest := fmt.Sprintf("%s.%s", target, prof.Name)
		hn, _ := os.Hostname()
		target := b.Start.Format("2006-01-02/15:04:05") + "-" + hn + "." + prof.Name
		p.uploader.UploadProfile(target, bytes.NewReader(prof.Data))
	}
	// p.uploader.Upload(bat, p.cfg.Service)
	return nil
}

func (p *profiler) outputDir(bat entity.Batch) error {
	if p.cfg.OutputDir == "" {
		return nil
	}
	// Basic ISO 8601 Format in UTC as the name for the directories.
	dir := bat.End.UTC().Format("20060102T150405Z")
	dirPath := filepath.Join(p.cfg.OutputDir, dir)
	// 0755 is what mkdir does, should be reasonable for the use cases here.
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	for _, prof := range bat.Profiles {
		filePath := filepath.Join(dirPath, prof.Name)
		// 0644 is what touch does, should be reasonable for the use cases here.
		if err := ioutil.WriteFile(filePath, prof.Data, 0644); err != nil {
			return err
		}
	}
	return nil
}

// interruptibleSleep sleeps for the given duration or until interrupted by the
// p.exit channel being closed.
func (p *profiler) interruptibleSleep(d time.Duration) {
	select {
	case <-p.exit:
	case <-time.After(d):
	}
}

// stop stops the profiler.
func (p *profiler) stop() {
	p.stopOnce.Do(func() {
		close(p.exit)
	})
	p.wg.Wait()
}

// WithPeriod specifies the interval at which to collect profiles.
func WithPeriod(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.Period = d
	}
}

// CPUDuration specifies the length at which to collect CPU profiles.
func CPUDuration(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.CpuDuration = d
	}
}

// WithProfileTypes specifies the profile types to be collected by the profiler.
func WithProfileTypes(types ...ProfileType) Option {
	return func(cfg *Config) {
		// reset the types and only use what the user has specified
		for k := range cfg.Types {
			delete(cfg.Types, k)
		}
		for _, t := range types {
			cfg.addProfileType(t)
		}
	}
}

// WithUploadTimeout specifies the timeout to use for uploading profiles. The
// default timeout is specified by DefaultUploadTimeout or the
// DD_PROFILING_UPLOAD_TIMEOUT env variable. Using a negative value or 0 will
// cause an error when starting the profiler.
func WithUploadTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.UploadTimeout = d
	}
}
