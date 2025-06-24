package versioning

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/mod/semver"
)

type Version struct {
	Prefix string
	Major  int64
	Minor  int64
	Patch  int64
}

func IsValid(version string, versionIsNotPrefixedByV bool) bool {
	prefixedVersion := version
	if versionIsNotPrefixedByV {
		prefixedVersion = fmt.Sprintf("v%s", version)
	}
	return semver.IsValid(prefixedVersion)
}

func Parse(version string, versionIsNotPrefixedByV bool) (*Version, error) {
	prefix := "v"
	if versionIsNotPrefixedByV {
		prefix = ""
	}

	if IsValid(version, versionIsNotPrefixedByV) {
		prefixedVersion := version
		if versionIsNotPrefixedByV {
			prefixedVersion = fmt.Sprintf("v%s", version)
		}
		majorSegment := semver.Major(prefixedVersion)
		majorMinorSegment := semver.MajorMinor(prefixedVersion)
		major, _ := strconv.ParseInt(majorSegment[1:], 10, 64)
		minor, _ := strconv.ParseInt(majorMinorSegment[len(majorSegment)+1:], 10, 64)
		patch, _ := strconv.ParseInt(prefixedVersion[len(majorMinorSegment)+1:], 10, 64)

		return &Version{
			Prefix: prefix,
			Major:  major,
			Minor:  minor,
			Patch:  patch,
		}, nil
	}
	return nil, errors.New("invalid version")
}

func (v *Version) String() string {
	return fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
}

func (v *Version) NextMajor() *Version {
	return &Version{
		Prefix: v.Prefix,
		Major:  v.Major + 1,
		Minor:  0,
		Patch:  0,
	}
}

func (v *Version) NextMinor() *Version {
	return &Version{
		Prefix: v.Prefix,
		Major:  v.Major,
		Minor:  v.Minor + 1,
		Patch:  0,
	}
}

func (v *Version) NextPatch() *Version {
	return &Version{
		Prefix: v.Prefix,
		Major:  v.Major,
		Minor:  v.Minor,
		Patch:  v.Patch + 1,
	}
}
