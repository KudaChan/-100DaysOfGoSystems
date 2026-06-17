package main

import (
	"reflect"
	"testing"
)

func TestReverseRange(t *testing.T) {
	t.Run("Even number of elements", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		ReverseSlice(s)
		expected := []int{5, 4, 3, 2, 1}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})

	t.Run("Odd number of elements", func(t *testing.T) {
		s := []string{"a", "b", "c"}
		ReverseSlice(s)
		expected := []string{"c", "b", "a"}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})
}

func TestReverseRange_(t *testing.T) {
	t.Run("Valid Range", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		err := ReverseRange(s, 1, 3)
		expected := []int{1, 4, 3, 2, 5}

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})

	t.Run("Out of bounds", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		err := ReverseRange(s, 1, 6)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Inverse Range", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		err := ReverseRange(s, 4, 2)
		if err == nil {
			t.Errorf("expected invalid range error, got nil")
		}
	})

	t.Run("Empty slice", func(t *testing.T) {
		s := []int{}
		err := ReverseRange(s, 0, 1)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
