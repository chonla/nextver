package versioning

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/mod/semver"
)

type Version struct {
	Major int64
	Minor int64
	Patch int64
}

func IsValid(version string) bool {
	return semver.IsValid(version)
}

func Parse(version string) (*Version, error) {
	if IsValid(version) {
		majorSegment := semver.Major(version)
		majorMinorSegment := semver.MajorMinor(version)
		major, _ := strconv.ParseInt(majorSegment[1:], 10, 64)
		minor, _ := strconv.ParseInt(majorMinorSegment[len(majorSegment)+1:], 10, 64)
		patch, _ := strconv.ParseInt(version[len(majorMinorSegment)+1:], 10, 64)

		return &Version{
			Major: major,
			Minor: minor,
			Patch: patch,
		}, nil
	}
	return nil, errors.New("invalid version")
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) NextMajor() *Version {
	return &Version{
		Major: v.Major + 1,
		Minor: 0,
		Patch: 0,
	}
}

func (v *Version) NextMinor() *Version {
	return &Version{
		Major: v.Major,
		Minor: v.Minor + 1,
		Patch: 0,
	}
}

func (v *Version) NextPatch() *Version {
	return &Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch + 1,
	}
}
