package versions

import (
	"reflect"
	"testing"
)

func TestVersionSorter(t *testing.T) {
	t.Parallel()
	versionStrings := []string{"10.1.0", "1.0.0", "1.10.0", "1.1.1", "2.3.1"}
	expectedSortedVersionStrings := []string{"1.0.0", "1.1.1", "1.10.0", "2.3.1", "10.1.0"}
	var versions []Version
	var sortedVersionStrings []string

	for _, s := range versionStrings {
		var v Version
		err := v.Parse(s)
		if err != nil {
			t.Error("Could not parse versions!")
		}
		versions = append(versions, v)
	}

	By(majorMinorBuild).Sort(versions)
	for _, v := range versions {
		sortedVersionStrings = append(sortedVersionStrings, v.ToString())
	}
	if !reflect.DeepEqual(sortedVersionStrings, expectedSortedVersionStrings) {
		t.Errorf("%v != %v", sortedVersionStrings, expectedSortedVersionStrings)
	}
}
