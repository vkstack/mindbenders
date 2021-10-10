package entity

import "time"

// batch is a collection of profiles of different types, collected at roughly the same time. It maps
// to what the dotpe UI calls a profile.
type Batch struct {
	Start, End time.Time
	Profiles   []*Profile
}

func (b *Batch) AddProfile(p *Profile) {
	b.Profiles = append(b.Profiles, p)
}
