package security

type Set[T comparable] map[T]bool

func (s *Set[T]) Add(value T) {
	(*s)[value] = true
}

func (s *Set[T]) Remove(value T) {
	delete(*s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	if ok {
		return true
	}
	return false
}

func SetFromSlice[T comparable](slice []T) Set[T] {
	set := make(Set[T])
	for _, value := range slice {
		set.Add(value)
	}
	return set
}
