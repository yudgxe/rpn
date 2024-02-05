package rpn

import "errors"

const EmptyStackTopIndex = -1

var ErrStackIsEmpty = errors.New("stack is empty")

type Stack[T any] struct {
	topIndex int
	buff     []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		buff:     make([]T, 0),
		topIndex: EmptyStackTopIndex,
	}
}

func (s *Stack[T]) Push(values ...T) {
	for _, v := range values {
		nextIndex := s.topIndex + 1
		if nextIndex < len(s.buff) {
			s.buff[nextIndex] = v
		} else {
			s.buff = append(s.buff, v)
		}
		s.topIndex = nextIndex
	}
}

func (s *Stack[T]) Pop() T {
	value, _ := s.PopWithError()
	return value
}

func (s *Stack[T]) PopWithError() (T, error) {
	if s.topIndex == EmptyStackTopIndex {
		var null T
		return null, ErrStackIsEmpty
	}
	v := s.buff[s.topIndex]
	s.topIndex--
	return v, nil
}

func (s *Stack[T]) PeekWithError() (T, error) {
	if s.topIndex == EmptyStackTopIndex {
		var null T
		return null, ErrStackIsEmpty
	}
	return s.buff[s.topIndex], nil
}

func (s *Stack[T]) Peek() T {
	value, _ := s.PeekWithError()
	return value
}

func (s *Stack[T]) Count() int {
	return s.topIndex + 1
}
