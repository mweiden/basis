package custom_alphabet

import (
  	"reflect"
	"testing"
)

func TestVersionSorter(t *testing.T) {
	t.Parallel()
	unsortedStrings := []string{"abc", "aba", "zzz", "zyz", "zyy", "aaa"}
	expectedStrings := []string{"zzz", "zyz", "zyy", "abc", "aba", "aaa"}

	runeIntMap := make(map[rune]int)
	for i, c := range "zyxwvutsrqponmlkjihgfedcba" {
		runeIntMap[rune(c)] = i
	}

	ordering := AlphabetOrdering{runeIntMap}

	By(ordering.customAlphabetical).Sort(unsortedStrings)

	if !reflect.DeepEqual(unsortedStrings, expectedStrings) {
		t.Errorf("%v != %v", unsortedStrings, expectedStrings)
	}
}
