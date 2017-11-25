package search

type Searchable interface {
	Len() int
	Compare(int, interface{}) int
}

func BinarySearch(s Searchable, point interface{}) (int, bool) {
	left := 0
	right := s.Len() - 1
	for left <= right {
		split := (left + right) / 2
		diff := s.Compare(split, point)
		if diff == 0 {
			return split, true
		} else if diff < 0 {
			left = split + 1
		} else {
			right = split - 1
		}
	}
	return left, false
}
