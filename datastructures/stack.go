package datastructures

import "errors"

var (
	EOS = errors.New("End of stack!")
)

type Stack struct {
	ary []interface{}
}

func (s *Stack) Top() interface{} {
	return s.ary[len(s.ary)-1]
}

func (s *Stack) Empty() bool {
	return len(s.ary) == 0
}

func (s *Stack) Push(i interface{}) {
	s.ary = append(s.ary, i)
}

func (s *Stack) Pop() (interface{}, error) {
	if s.Empty() {
		return nil, EOS
	} else {
		result := s.ary[len(s.ary)-1]
		s.ary = s.ary[:len(s.ary)-1]
		return result, nil
	}
}
