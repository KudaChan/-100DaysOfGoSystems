# I Built a Generic Dynamic Array From Scratch in Go, and It Taught Me How the Runtime Actually Manages Memory

### Why building standard data structures from first principles is the ultimate cheat code for understanding Go’s internals, low-level pointers, and the garbage collector.

When writing modern backend applications in Go, we rarely think twice about using slices. We type `append(slice, item)` and let the Go runtime manage the rest. It feels clean, seamless, and magically performant.

But magic is a liability in production systems.

As part of my **#100DaysOfGoSystems** engineering challenge, I decided to pull back the curtain. My objective for Day 1 was simple yet deceptive: **Implement a fully generic Dynamic Array (Vector) from scratch without using Go’s built-in `append()` or `copy()` primitives.** I wanted to face raw memory management, manual sizing mechanics, and the actual mathematical trade-offs behind geometric scaling.

What I expected was a straightforward exercise in structural loops. What I actually found was a masterclass in pointer receiver mechanics, slice header copying, and a silent, high-severity memory leak trap built into the Go Garbage Collector.

Here is the architectural post-mortem of what I built, what broke, and how the Go runtime handles memory under the hood.

---

## The Core Blueprint: Replicating the Engine

To build a dynamic array, you have to realize that primitive arrays in Go are fixed-size blocks of memory whose length is determined strictly at compile-time. To make an array dynamic, you must encapsulate it within a management layer that handles resizing allocations on the fly.

I defined our generic structure like this:

```go
type DynamicArray[T any] struct {
    data     []T // The underlying contiguous block of memory
    size     int // Number of elements currently occupied
    capacity int // Total memory slots allocated
}

```

This structure mimics what Go natively tracks inside a slice header. To mimic a low-level engine, the API needed four foundational components:

1. `Append(element T)`: O(1) average-time insertion.
2. `Get(index int) (T, error)`: O(1) random memory access with defensive boundary protection.
3. `Delete(index int) error)`: Contiguous item removal by left-shifting indices.
4. `grow()` / `shrink()`: Internal memory reallocation engines.

---

## Deep Dive 1: The Geometry of $O(1)$ and Preventing "Thrashing"

If you resize an underlying array by adding a fixed number of slots (e.g., `capacity + 10`) every time it fills up, your `Append` operations will slowly degrade into a performance nightmare. At a size of $N$, every single insert will trigger a costly $O(N)$ allocation and data copy loop.

To achieve an **Amortized Time Complexity of $O(1)$**, the internal `grow()` engine must scale *geometrically*. When `size == capacity`, the array doubles its size ($2\times$). By growing exponentially, the expensive $O(N)$ copy routine happens less and less frequently, meaning that over millions of inserts, the average cost per insertion remains virtually instantaneous.

```go
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

```

### The Inverted Problem: Shrinking without Thrashing

Many engineers remember to grow their arrays, but they forget to down-size them when elements are deleted, leading to long-running memory bloat in daemon microservices.

However, if you shrink the array by half the moment it falls below $50\%$ capacity, you create a system vulnerability known as **thrashing**. Imagine an array hovering right at a capacity threshold: alternating calls of `Append` and `Delete` will force your service to continuously allocate and deallocate memory on every single operation.

To protect system CPU, I engineered an **asymmetric shrinking strategy**: the backing array only down-sizes by half ($0.5\times$) when its load factor drops to or below **$25\%$ ($1/4$ capacity)**. This creates a buffer zone that prevents constant allocation cycles.

---

## Deep Dive 2: The Two Memory Traps I Had to Solve

While implementing the code, I ran directly into two classic low-level Go pitfalls.

### Trap #1: Struct Value Duplication on the Stack

In Go, absolutely everything is passed by value. When you define a struct method, you have to choose between a value receiver `func (da DynamicArray[T])` and a pointer receiver `func (da *DynamicArray[T])`.

If I had used a value receiver for `Append`, the Go runtime would have duplicated the `DynamicArray` header onto the stack when the method was called. The code would execute, increment `size`, swap the internal pointer during `grow()`, and then... *completely discard those changes* the millisecond the function ended. The caller’s original array would remain entirely unmodified.

Using **pointer receivers** ensures that we are mutating the caller's actual memory coordinates directly.

### Trap #2: The Silent Zero-Value Reference Leak

This was the most sinister bug of all. In the `Delete` method, when an element is removed, all items to its right are shifted one index to the left:

```go
// Shifting items left
for i := index; i < da.size-1; i++ {
    da.data[i] = da.data[i+1]
}
da.size--

```

While this correctly decrements the publicly accessible `size`, look at what is happening under the hood. The very last index slot of the array *still contains a duplicate reference to the item that was shifted*.

If your dynamic array is storing large structs or memory pointers, **the Go Garbage Collector (GC) will see that stale trailing reference as an active object link.** Because it is still technically referenced in the array's backing block, the GC will never free that memory, resulting in a severe, slow-burning memory leak.

The fix requires systems-level discipline: we must explicitly clear out the vacated slot using the generic type's **zero value** to break the reference chain:

```go
var zero T
da.data[da.size] = zero // Breaks the reference chain, allowing immediate GC cleanup!

```

---

## Production-Ready Verification via Table-Driven Tests

In a production systems environment, code without automated verification is broken by design. To prove this custom structure was bulletproof against runtime panics and memory leaks, I implemented an idiomatic, table-driven unit test suite.

The test suite explicitly validated three crucial edge cases:

1. **Zero-Capacity Assertions:** Ensuring that instantiating an array with an initial capacity of `0` smoothly transitions to a capacity of `1` on the first append without panicking out-of-bounds.
2. **Resizing Accuracy:** Confirming that capacity precisely doubles and caps off correctly across data boundaries.
3. **GC Integrity Verification:** Introspecting the underlying slice slot directly after a deletion to verify that the trailing data slot was successfully zeroed out.

```go
t.Run("Delete and Memory Leak Prevention", func(t *testing.T) {
    da := NewDynamicArray[int](4)
    da.Append(10)
    da.Append(20)
    da.Append(30)

    da.Delete(1) // Remove 20

    // Structural Check: Enforce that the trailing slot is completely zeroed out
    if da.data[da.size] != 0 {
        t.Errorf("GC Trap! Stale memory index was not zeroed out")
    }
})

```

---

## Summary of Architectural Takeaways

Stripping away Go's built-in abstractions on Day 1 completely reframed how I visualize memory allocation in high-throughput backends:

1. **Memory Ownership Matters:** Just because an item is inaccessible via your API doesn't mean it's invisible to the Garbage Collector. Always clean up references when building custom collection types.
2. **Allocation Trashing is Real:** Asymmetric allocation profiles ($2\times$ scaling vs $1/4$ capacity shrinking) are vital for keeping CPU cycles predictable under unstable, heavy workloads.
3. **Embrace Generics Safely:** Go Generics (`[T any]`) deliver flawless type safety without the heavy performance taxes or runtime interface type-assertions of the past, provided you manage fallback declarations using clean `var zero T` zero-value idioms.

Day 1 is officially in the books, and the foundational mental model is set. Up next for Day 2, I am tackling **In-Place Slice Reversals** to dive deep into memory swaps and index optimizations.

*The complete, verified source code for this challenge is available on my GitHub repository. Follow along as I build out the next 99 challenges!*

---

**#GoLang #SystemsEngineering #BackendArchitecture #SoftwareDevelopment #100DaysOfCode**
