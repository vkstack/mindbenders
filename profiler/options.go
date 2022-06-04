package profiler

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.com/dotpe/mindbenders/profiler/uploader"
)

const (
	// DefaultPeriod specifies the default period at which profiles will be collected.
	DefaultPeriod = 30 * time.Second

	// DefaultDuration specifies the default length of the CPU profile snapshot.
	DefaultDuration = time.Second * 15

	// DefaultUploadTimeout specifies the default timeout for uploading profiles.
	DefaultUploadTimeout = 10 * time.Second
)

var defaultUploader = uploader.NewNullUploader()

type Config struct {
	// targetURL is the upload destination URL. It will be set by the profiler on start to either apiURL or agentURL
	// based on the other options.
	Service, Host     string
	Types             map[ProfileType]struct{}
	Period            time.Duration
	CpuDuration       time.Duration
	UploadTimeout     time.Duration
	MaxGoroutinesWait int
	OutputDirFn       func(time.Time) string
	uploader          uploader.IProfileUploader
}

func (c *Config) AddProfileType(t ProfileType) {
	if c.Types == nil {
		c.Types = make(map[ProfileType]struct{})
	}
	c.Types[t] = struct{}{}
}

var defaultProfileTypes = []ProfileType{CPUProfile}

func (c *Config) addProfileType(t ProfileType) {
	if c.Types == nil {
		c.Types = make(map[ProfileType]struct{})
	}
	c.Types[t] = struct{}{}
}

func defaultConfig() (*Config, error) {
	hn, _ := os.Hostname()
	c := Config{
		Host:              hn,
		Service:           filepath.Base(os.Args[0]),
		Period:            DefaultPeriod,
		CpuDuration:       DefaultDuration,
		UploadTimeout:     DefaultUploadTimeout,
		MaxGoroutinesWait: 1000, // arbitrary value, should limit STW to ~30ms
		uploader:          defaultUploader,
	}
	for _, t := range defaultProfileTypes {
		c.addProfileType(t)
	}
	if v := os.Getenv("MINDBENDERS_SERVICE"); v != "" {
		WithService(v)(&c)
	}
	if v := os.Getenv("PROFILING_WAIT_PROFILE_MAX_GOROUTINES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("PROFILING_WAIT_PROFILE_MAX_GOROUTINES: %s", err)
		}
		c.MaxGoroutinesWait = n
	}
	return &c, nil
}

// An Option is used to configure the profiler's behaviour.
type Option func(*Config)

// WithService specifies the service name to attach to a profile.
func WithService(name string) Option {
	return func(cfg *Config) {
		cfg.Service = name
	}
}

func WithUploader(uploader uploader.IProfileUploader) Option {
	return func(cfg *Config) {
		cfg.uploader = uploader
	}
}

/*
	The below targetsetter is a very general pathsetter
		the location of the profiles will be as follows
	--------{ENV}	>>	{Service}	>>	{Date (YYYY-MM-DD)}
		>>	{Host}	>>	{Hour}	>>	{Minute:Second}	>>	{{{all profiles}}}
*/

func WithTargetSetter(cfg *Config) {
	cfg.OutputDirFn = func(t time.Time) string {
		return path.Join(os.Getenv("ENV"), cfg.Service, t.Format("2006-01-02"), cfg.Host, t.Format("15/04:05.000"))
	}
}
