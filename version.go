package main

import (
	"errors"
	"fmt"
)

var ErrInvalidVersion = errors.New("invalid version")

type version struct {
	Minor     uint8
	Patch     uint8
	Candidate uint8
	Beta      uint8
}

func parse(s string) (version, error) {
	var v version
	err := v.UnmarshalText([]byte(s))

	return v, err
}

func (v version) String() string {
	if v == (version{}) {
		return "gotip"
	}

	var pbc string

	if v.Patch > 0 {
		pbc = fmt.Sprintf(".%d", v.Patch)
	}

	if v.Beta > 0 {
		pbc += fmt.Sprintf("beta%d", v.Beta)
	} else if v.Candidate > 0 {
		pbc += fmt.Sprintf("rc%d", v.Candidate)
	}

	return fmt.Sprintf("go1.%d%s", v.Minor, pbc)
}

func (v version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *version) UnmarshalText(data []byte) error {
	s := string(data)

	if s == "gotip" {
		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%d.%d", &v.Minor, &v.Patch); err == nil {
		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%dbeta%d", &v.Minor, &v.Beta); err == nil {
		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%drc%d", &v.Minor, &v.Candidate); err == nil {
		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%d", &v.Minor); err == nil {
		return nil
	}

	return ErrInvalidVersion
}

type versions []version

func (v versions) Len() int {
	return len(v)
}

func (v versions) Less(i, j int) bool {
	if v[i] == (version{}) {
		return false
	}

	if v[j] == (version{}) {
		return true
	}

	if v[i].Minor != v[j].Minor {
		return v[i].Minor < v[j].Minor
	}

	if v[i].Patch != v[j].Patch {
		return v[i].Patch < v[j].Patch
	}

	if v[i].Candidate != v[j].Candidate {
		return v[i].Candidate < v[j].Candidate
	}

	return v[i].Beta < v[j].Beta
}

func (v versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
