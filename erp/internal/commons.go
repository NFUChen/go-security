package internal

// Define a set using a map
type Set[T comparable] map[T]bool

// Add an element to the set
func (s *Set[T]) Add(value T) {
	(*s)[value] = true
}

// Remove an element from the set
func (s *Set[T]) Remove(value T) {
	delete(*s, value)
}

// Check if an element exists in the set
func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	if ok {
		return true
	}
	return false
}
