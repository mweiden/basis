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
	}
	result := s.ary[len(s.ary)-1]
	s.ary = s.ary[:len(s.ary)-1]
	return result, nil
}

func (s *Stack) Peek() (interface{}, error) {
	if s.Empty() {
		return nil, EOS
	}
	return s.ary[len(s.ary)-1], nil
}
