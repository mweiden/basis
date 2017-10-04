package phone_numbers

import (
	"fmt"
	"strconv"
)

var (
	LETTER_TO_INT = map[rune]rune{
		'A': '2',
		'B': '2',
		'C': '2',
		'D': '3',
		'E': '3',
		'F': '3',
		'G': '4',
		'H': '4',
		'I': '4',
		'J': '5',
		'K': '5',
		'L': '5',
		'M': '6',
		'N': '6',
		'O': '6',
		'P': '7',
		'Q': '7',
		'R': '7',
		'S': '7',
		'T': '8',
		'U': '8',
		'V': '8',
		'W': '9',
		'X': '9',
		'Y': '9',
		'Z': '9',
	}
)

type Dictionary struct {
	words []string
}

func (d Dictionary) ToPhoneNumberMap() func(string) []string {
	m := make(map[string][]string)
	for _, w := range d.words {
		key := ""
		for _, c := range w {
			key += string(LETTER_TO_INT[c])
		}
		m[key] = append(m[key], w)
	}
	return func(key string) []string {
		return m[key]
	}
}

type With func(intsStr string) []string

func key(ints []int) string {
	key := ""
	for _, i := range ints {
		key += strconv.Itoa(i)
	}
	return key
}

func (w With) PossibleWords(ints []int) []string {
	return w(key(ints))
}

func (w With) PossibleWordsForPhoneNumber(phoneNumber []int) []string {
	var possible []string
	country := strconv.Itoa(phoneNumber[0])
	var area string
	for _, n := range phoneNumber[1:4] {
		area += strconv.Itoa(n)
	}
	for _, w1 := range w.PossibleWords(phoneNumber[4:7]) {
		for _, w2 := range w.PossibleWords(phoneNumber[7:11]) {
			ns := fmt.Sprintf("%s-%s-%s-%s", country, area, w1, w2)
			possible = append(possible, ns)
		}
	}
	return possible
}
