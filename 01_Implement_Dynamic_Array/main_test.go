package main

import (
	"errors"
	"testing"
)

func TestDynamicArray(t *testing.T) {
	t.Run("Append and Grow Mechanism", func(t *testing.T) {
		da := NewDynamicArray[int](2)

		da.Append(10)
		da.Append(20)

		if da.capacity != 2 {
			t.Errorf("expected cap 2, got %d", da.capacity)
		}

		da.Append(30)
		if da.capacity != 4 {
			t.Errorf("expected cap to double to 4, got %d", da.capacity)
		}

		if da.size != 3 {
			t.Errorf("expected size 3, got %d", da.size)
		}
	})

	t.Run("Get and Bounds Checking", func(t *testing.T) {
		da := NewDynamicArray[string](5)
		da.Append("Go")

		val, err := da.Get(0)
		if err != nil || val != "Go" {
			t.Errorf("expected 'Go' with no error, got %v, %v", val, err)
		}

		_, err = da.Get(1)
		if !errors.Is(err, ErrorOutOfBound) {
			t.Errorf("expected ErrorOutOfBound, got %v", err)
		}
	})

	t.Run("Delete and Memory Leak Prevention", func(t *testing.T) {
		da := NewDynamicArray[int](4)
		da.Append(10)
		da.Append(20)
		da.Append(30)

		newSize, err := da.Delete(1)
		if err != nil || newSize != 2 {
			t.Errorf("expected new size 2, got %d", newSize)
		}

		val, _ := da.Get(1)
		if val != 30 {
			t.Errorf("expected element at index 1 to shift to 30, got %v", val)
		}

		if da.data[da.size] != 0 {
			t.Error("GC Trap! Stale memory index was not zeroed out")
		}
	})
}
