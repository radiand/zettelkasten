package testutils

// Cycle allows to iterate slice indefinitely.
type Cycle[T any] struct {
	array []T
	last  int
}

// Next returns consecutive elements from the slice. If slice is exhausted, it
// starts over.
func (self *Cycle[T]) Next() T {
	if len(self.array) > self.last {
		current := self.last
		self.last++
		return self.array[current]
	}
	self.last = 0
	return self.array[0]
}

// Enqueue adds new element to iterate over.
func (self *Cycle[T]) Enqueue(element ...T) {
	self.array = append(self.array, element...)
}

// NewCycle creates new Cycle instance.
func NewCycle[T any](items ...T) Cycle[T] {
	return Cycle[T]{
		array: items,
		last:  0,
	}
}
