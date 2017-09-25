package custom_alphabet

import "sort"

type AlphabetOrdering struct {
	runeIntMap map[rune]int
}

func (a AlphabetOrdering) Validate() bool {
	counts := make(map[int]int)

	for _, i := range a.runeIntMap {
		_, ok := counts[i]
		if !ok {
			counts[i] = 1
		} else {
			return false
		}
	}
	return true
}

type customAlphabetSorter struct {
	strings []string
	runeIntMap map[rune]int
	by      func(p1 string, p2 string) bool
}

type By func(p1 string, p2 string) bool

func (by By) Sort(strings []string) {
	cas := &customAlphabetSorter{
		strings: strings,
		by:      by,
	}
	sort.Sort(cas)
}

func (cas *customAlphabetSorter) Swap(i, j int) {
	cas.strings[i], cas.strings[j] = cas.strings[j], cas.strings[i]
}

func (cas *customAlphabetSorter) Len() int {
	return len(cas.strings)
}

func (cas *customAlphabetSorter) Less(i, j int) bool {
	return cas.by(cas.strings[i], cas.strings[j])
}

func (o AlphabetOrdering) customAlphabetical(s1 string, s2 string) bool {
	diff := 0
	for i := 0; i < min(len(s1), len(s2)); i++ {
		diff = o.runeIntMap[rune(s2[i])] - o.runeIntMap[rune(s1[i])]
		if diff != 0 {
			return diff > 0
		}
	}
	diff = len(s2) - len(s1)
	return diff > 0
}

func min(i1 int, i2 int) int {
	if i1 < i2 {
		return i1
	} else {
		return i2
	}
}
