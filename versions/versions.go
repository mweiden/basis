package versions

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	versionNumbers []int
}

func (v *Version) Parse(s string) error {
	versionStrings := strings.Split(s, ".")
	if len(versionStrings) > 3 {
		return errors.New("Not a valid version string!")
	}
	for _, str := range versionStrings {
		versionInt, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		v.versionNumbers = append(v.versionNumbers, versionInt)
	}
	return nil
}

func (v *Version) ToString() string {
	versionNumberStrings := make([]string, len(v.versionNumbers))
	for i, v := range v.versionNumbers {
		versionNumberStrings[i] = strconv.Itoa(v)
	}
	return strings.Join(versionNumberStrings, ".")
}

type versionSorter struct {
	versions []Version
	by       func(p1, p2 *Version) bool // Closure used in the Less method.
}

type By func(p1, p2 *Version) bool

func (by By) Sort(versions []Version) {
	vs := &versionSorter{
		versions: versions,
		by:       by,
	}
	sort.Sort(vs)
}

func (vs *versionSorter) Swap(i, j int) {
	vs.versions[i], vs.versions[j] = vs.versions[j], vs.versions[i]
}

func (vs *versionSorter) Len() int {
	return len(vs.versions)
}

func (s *versionSorter) Less(i, j int) bool {
	return s.by(&s.versions[i], &s.versions[j])
}

func majorMinorBuild(v1, v2 *Version) bool {
	for i := 0; i < len(v1.versionNumbers); i++ {
		diff := v1.versionNumbers[i] - v2.versionNumbers[i]
		if diff < 0 {
			return true
		} else if diff > 0 {
			return false
		}
	}
	return false
}
