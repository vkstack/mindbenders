package profiler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime/pprof"
	"time"

	"gitlab.com/dotpe/mindbenders/profiler/entity"
	pprofutils "gitlab.com/dotpe/mindbenders/profiler/internal/pprofutils"

	pprofile "github.com/google/pprof/profile"
)

// ProfileType represents a type of profile that the profiler is able to run.
type ProfileType int

const (
	// HeapProfile reports memory allocation samples; used to monitor current
	// and historical memory usage, and to check for memory leaks.
	HeapProfile ProfileType = iota
	// CPUProfile determines where a program spends its time while actively consuming
	// CPU cycles (as opposed to while sleeping or waiting for I/O).
	CPUProfile
	// GoroutineProfile reports stack traces of all current goroutines
	GoroutineProfile
	// MetricsProfile reports top-line metrics associated with user-specified profiles
	MetricsProfile
)

// profileTypes maps every ProfileType to its implementation.
var profileTypes = map[ProfileType]profileType{
	CPUProfile: {
		Name:     "cpu",
		Filename: "cpu.pprof",
		Collect: func(_ profileType, p *profiler) ([]byte, error) {
			var buf bytes.Buffer
			if err := startCPUProfile(&buf); err != nil {
				return nil, err
			}
			p.interruptibleSleep(p.cfg.CpuDuration)
			stopCPUProfile()
			return buf.Bytes(), nil
		},
	},
	GoroutineProfile: {
		Name:     "goroutine",
		Filename: "goroutines.pprof",
		Collect:  collectGenericProfile,
	},
	// HeapProfile is complex due to how the Go runtime exposes it. It contains 4
	// sample types alloc_objects/count, alloc_space/bytes, inuse_objects/count,
	// inuse_space/bytes. The first two represent allocations over the lifetime
	// of the process, so we do delta profiling for them. The last two are
	// snapshots of the current heap state, so we leave them as-is.
	HeapProfile: {
		Name:     "heap",
		Filename: "heap.pprof",
		Delta: &pprofutils.Delta{SampleTypes: []pprofutils.ValueType{
			{Type: "alloc_objects", Unit: "count"},
			{Type: "alloc_space", Unit: "bytes"},
		}},
		Collect: collectGenericProfile,
	},
	MetricsProfile: {
		Name:     "metrics",
		Filename: "metrics.json",
		Collect: func(_ profileType, p *profiler) ([]byte, error) {
			var buf bytes.Buffer
			err := p.met.report(now(), &buf)
			return buf.Bytes(), err
		},
	},
}

// profileType holds the implementation details of a ProfileType.
type profileType struct {
	// Type gets populated automatically by ProfileType.lookup().
	Type ProfileType
	// Name specifies the profile name as used with pprof.Lookup(name) (in
	// collectGenericProfile) and returned by ProfileType.String(). For profile
	// types that don't use this approach (e.g. CPU) the name isn't used for
	// anything.
	Name string
	// Filename is the filename used for uploading the profile to the dotpe
	// backend which is aware of them. Delta profiles are prefixed with "delta-"
	// automatically. In theory this could be derrived from the Name field, but
	// this isn't done due to idiosyncratic filename used by the
	// GoroutineProfile.
	Filename string
	// Delta controls if this profile should be generated as a delta profile.
	// This is useful for profiles that represent samples collected over the
	// lifetime of the process (i.e. heap, block, mutex). If nil, no delta
	// profile is generated.
	Delta *pprofutils.Delta
	// Collect collects the given profile and returns the data for it. Most
	// profiles will be in pprof format, i.e. gzip compressed proto buf data.
	Collect func(profileType, *profiler) ([]byte, error)
}

// lookup returns t's profileType implementation.
func (t ProfileType) lookup() profileType {
	c, ok := profileTypes[t]
	if ok {
		c.Type = t
		return c
	}
	return profileType{
		Type:     t,
		Name:     "unknown",
		Filename: "unknown",
		Collect: func(_ profileType, _ *profiler) ([]byte, error) {
			return nil, errors.New("profile type not implemented")
		},
	}
}

func collectGenericProfile(t profileType, _ *profiler) ([]byte, error) {
	var buf bytes.Buffer
	err := lookupProfile(t.Name, &buf, 0)
	return buf.Bytes(), err
}

// String returns the name of the profile.
func (t ProfileType) String() string {
	return t.lookup().Name
}

// Filename is the identifier used on upload.
func (t ProfileType) Filename() string {
	return t.lookup().Filename
}

// Tag used on profile metadata
func (t ProfileType) Tag() string {
	return fmt.Sprintf("profile_type:%s", t)
}

func (p *profiler) runProfile(pt ProfileType) ([]*entity.Profile, error) {
	t := pt.lookup()
	// Collect the original profile as-is.
	data, err := t.Collect(t, p)
	if err != nil {
		return nil, err
	}
	profs := []*entity.Profile{{
		Name: t.Filename,
		Data: data,
	}}
	// Compute the deltaProf (will be nil if not enabled for this profile type).
	deltaProf, err := p.deltaProfile(t, data)
	if err != nil {
		return nil, fmt.Errorf("delta profile error: %s", err)
	}
	// Report metrics and append deltaProf if not nil.
	if deltaProf != nil {
		profs = append(profs, deltaProf)
	}
	return profs, nil
}

// deltaProfile derives the delta profile between curData and the previous
// profile. For profile types that don't have delta profiling enabled, it
// simply returns nil, nil.
func (p *profiler) deltaProfile(t profileType, curData []byte) (*entity.Profile, error) {
	// Not all profile types use delta profiling, return nil if this one doesn't.
	if t.Delta == nil {
		return nil, nil
	}
	curProf, err := pprofile.ParseData(curData)
	if err != nil {
		return nil, fmt.Errorf("delta prof parse: %v", err)
	}
	var deltaData []byte
	if prevProf := p.prev[t.Type]; prevProf == nil {
		// First time deltaProfile gets called for a type, there is no prevProf. In
		// this case we emit the current profile as a delta profile.
		deltaData = curData
	} else {
		// Delta profiling is also implemented in the Go core, see commit below.
		// Unfortunately the core implementation isn't resuable via a API, so we do
		// our own delta calculation below.
		// https://github.com/golang/go/commit/2ff1e3ebf5de77325c0e96a6c2a229656fc7be50#diff-94594f8f13448da956b02997e50ca5a156b65085993e23bbfdda222da6508258R303-R304
		deltaProf, err := t.Delta.Convert(prevProf, curProf)
		if err != nil {
			return nil, fmt.Errorf("delta prof merge: %v", err)
		}
		// TimeNanos is supposed to be the time the profile was collected, see
		// https://github.com/google/pprof/blob/master/proto/profile.proto.
		deltaProf.TimeNanos = curProf.TimeNanos
		// DurationNanos is the time period covered by the profile.
		deltaProf.DurationNanos = curProf.TimeNanos - prevProf.TimeNanos
		deltaBuf := &bytes.Buffer{}
		if err := deltaProf.Write(deltaBuf); err != nil {
			return nil, fmt.Errorf("delta prof write: %v", err)
		}
		deltaData = deltaBuf.Bytes()
	}
	// Keep the most recent profiles in memory for future diffing. This needs to
	// be taken into account when enforcing memory limits going forward.
	p.prev[t.Type] = curProf
	return &entity.Profile{
		Name: "delta-" + t.Filename,
		Data: deltaData,
	}, nil
}

var (
	// startCPUProfile starts the CPU profile; replaced in tests
	startCPUProfile = pprof.StartCPUProfile
	// stopCPUProfile stops the CPU profile; replaced in tests
	stopCPUProfile = pprof.StopCPUProfile
)

// lookpupProfile looks up the profile with the given name and writes it to w. It returns
// any errors encountered in the process. It is replaced in tests.
var lookupProfile = func(name string, w io.Writer, debug int) error {
	prof := pprof.Lookup(name)
	if prof == nil {
		return errors.New("profile not found")
	}
	return prof.WriteTo(w, debug)
}

// now returns current time in UTC.
func now() time.Time {
	return time.Now().UTC()
}
