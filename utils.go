package main

import (
	"fmt"
	"time"
)

type Stack struct {
	size       int32
	topElement *StackElement
}

type StackElement struct {
	//interface support any type
	value interface{}
	next  *StackElement
}

func (s *Stack) Size() int32 {
	return s.size
}

func (s *Stack) Push(element interface{}) {
	newElement := StackElement{
		value: element,
		next:  s.topElement,
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

func getCurrentTs() int64 {
	return time.Now().UnixNano() / 1000000
}

func getFundsAsString(amount float64) string {
	if amount == 0 {
		return ""
	}
	return fmt.Sprintf("%.2f", float64(amount))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
