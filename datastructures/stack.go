package main

import "errors"

var (
	EOS = errors.New("End of stack!")
)

type Stack struct {
	ary []int
}

func (s *Stack) Top() int {
	return s.ary[len(s.ary)-1]
}

func (s *Stack) Empty() bool {
	return len(s.ary) == 0
}

func (s *Stack) Push(i int) {
	s.ary = append(s.ary, i)
}

func (s *Stack) Pop() (error, int) {
	if s.Empty() {
		return EOS, -1
	} else {
		result := s.ary[len(s.ary)-1]
		s.ary = s.ary[:len(s.ary)-1]
		return nil, result
	}
}
