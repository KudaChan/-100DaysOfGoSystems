package main

import (
	"fmt"
)

// ReverseSlice reverses any generic slice in-place with O(1) extra space
func ReverseSlice[T any](s []T) {
	r := len(s) - 1
	l := 0

	for l < r {
		s[l], s[r] = s[r], s[l]
		l++
		r--
	}
}

// Challenge Extension: ReverseRange reverses only a specific window of a slice in-place
func ReverseRange[T any](s []T, start, end int) error {
	if start < 0 || end >= len(s) || start > end {
		return fmt.Errorf("start: %d, end: %d are out of bounds", start, end)
	}

	ReverseSlice(s[start : end+1])
	return nil
}

func main() {
	// Sanity Check 1: Full Reversal
	nums := []int{10, 20, 30, 40, 50}
	ReverseSlice(nums)
	fmt.Println("Reversed:", nums) // Expected: [50, 40, 30, 20, 10]

	// Sanity Check 2: Range Reversal
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	// Reverse only "banana", "cherry", "date" (indices 1 to 3)
	ReverseRange(words, 1, 3)
	fmt.Println("Range Reversed:", words) // Expected: ["apple", "date", "cherry", "banana", "elderberry"]
}
