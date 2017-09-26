package text_validator

import (
	"github.com/basis/datastructures"
)

var (
	OPENINGS              = []rune{'(', '{', '['}
	CLOSINGS              = []rune{')', '}', ']'}
	PARENS_BRACES_MAPPING = map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
)

type ParensValidator struct {
	state datastructures.Stack
}

func (pv *ParensValidator) Validate(next rune) bool {
	if contains(CLOSINGS, next) {
		shouldBe := PARENS_BRACES_MAPPING[next]
		if !pv.state.Empty() && pv.state.Top() == shouldBe {
			pv.state.Pop()
			return true
		} else {
			return false
		}
	}
	return true
}

func (pv *ParensValidator) Observe(next rune) {
	if contains(OPENINGS, next) {
		pv.state.Push(next)
	}
}

func contains(set []rune, ele rune) bool {
	for _, val := range set {
		if ele == val {
			return true
		}
	}
	return false
}
