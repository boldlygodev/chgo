package main

import (
	"errors"
	"fmt"
)

const (
	stageRelease uint8 = 255 - iota
	stageCandidate
	stageBeta
)

var ErrInvalidVersion = errors.New("invalid version")

type version struct {
	Minor      uint8
	Patch      uint8
	Stage      uint8
	Prerelease uint8
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

	if v.Stage == stageRelease && v.Patch > 0 {
		pbc = fmt.Sprintf(".%d", v.Patch)
	}

	if v.Stage == stageCandidate {
		pbc += fmt.Sprintf("rc%d", v.Prerelease)
	}

	if v.Stage == stageBeta {
		pbc += fmt.Sprintf("beta%d", v.Prerelease)
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
		v.Stage = stageRelease

		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%drc%d", &v.Minor, &v.Prerelease); err == nil {
		v.Stage = stageCandidate

		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%dbeta%d", &v.Minor, &v.Prerelease); err == nil {
		v.Stage = stageBeta

		return nil
	}

	if _, err := fmt.Sscanf(s, "go1.%d", &v.Minor); err == nil {
		v.Stage = stageRelease

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

	if v[i].Stage != v[j].Stage {
		return v[i].Stage < v[j].Stage
	}

	return v[i].Prerelease < v[j].Prerelease
}

func (v versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
