package main

type Stack struct {
	size int32
	topElement *StackElement
}

type StackElement struct {
	//interface supports any type
	value interface {}
	next *StackElement
}

func (s *Stack) Size() int32 {
	return s.size
}

func (s *Stack) Push(element interface{}) {
	newElement := StackElement{ 
		value: element,
		next: s.topElement,
	}
	s.topElement = &newElement
	s.size++
}

func (s *Stack) Pop() interface{} {
	if s.size > 0 {
		value := s.topElement.value
		s.topElement = s.topElement.next
		s.size--
		return value
	}
	return nil
}