package main

import (
	"errors"
	"fmt"
)

// ErrorOutOfBound is returned when the index is out of bound.
var ErrorOutOfBound = errors.New("index out of bound")

// DynamicArray is a generic dynamic array implementation.
type DynamicArray[T any] struct {
	data     []T
	size     int
	capacity int
}

// NewDynamicArray creates a new DynamicArray with the given initial capacity.
func NewDynamicArray[T any](initialCapacity int) *DynamicArray[T] {
	if initialCapacity < 0 {
		initialCapacity = 0
	}
	return &DynamicArray[T]{
		data:     make([]T, initialCapacity),
		size:     0,
		capacity: initialCapacity,
	}
}

// Append adds an element to the end of the DynamicArray.
func (da *DynamicArray[T]) Append(element T) {
	if da.size == da.capacity {
		da.grow()
	}
	da.data[da.size] = element
	da.size++
}

// Get returns the element at the given index.
func (da *DynamicArray[T]) Get(index int) (T, error) {
	if index < 0 || index >= da.size {
		var zero T
		return zero, ErrorOutOfBound
	}
	return da.data[index], nil
}

// Delete removes the element at the given index.
func (da *DynamicArray[T]) Delete(index int) (int, error) {
	if index < 0 || index >= da.size {
		return 0, ErrorOutOfBound
	}
	for i := index; i < da.size-1; i++ {
		da.data[i] = da.data[i+1]
	}

	da.size--
	var zero T
	da.data[da.size] = zero

	if da.capacity > 4 && da.size <= da.capacity/4 {
		da.shrink()
	}

	return da.size, nil
}

// grow increases the capacity of the DynamicArray.
func (da *DynamicArray[T]) grow() {
	if da.capacity == 0 {
		da.capacity = 1
	} else {
		da.capacity *= 2
	}

	newSlice := make([]T, da.capacity)
	for i := 0; i < da.size; i++ {
		newSlice[i] = da.data[i]
	}
	da.data = newSlice
}

// shrink decreases the capacity of the DynamicArray.
func (da *DynamicArray[T]) shrink() {
	da.capacity /= 2
	newSlice := make([]T, da.capacity)
	for i := 0; i < da.size; i++ {
		newSlice[i] = da.data[i]
	}
	da.data = newSlice
}

func main() {
	da := NewDynamicArray[int](2)
	da.Append(10)
	da.Append(20)
	da.Append(30)

	fmt.Printf("Before Delete - Size: %d, Cap: %d, Data: %v\n", da.size, da.capacity, da.data)

	da.Delete(1)

	fmt.Printf("After Delete  - Size: %d, Cap: %d, Data: %v\n", da.size, da.capacity, da.data)
}
