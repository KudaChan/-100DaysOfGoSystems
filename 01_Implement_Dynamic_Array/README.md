# Day 1: Generic Dynamic Array from Scratch

## Architectural Overview

This project implements a fully generic, production-grade dynamic array (`DynamicArray[T]`) from scratch in Go. Since primitive arrays in Go have fixed boundaries fixed at compile-time, this structure abstracts low-level memory allocations to simulate a contiguous, dynamically scaling vector.

To balance memory footprint and CPU efficiency, the implementation utilizes an **Amortized Time Complexity Strategy** modeled after systems-level vector engines:

* **Geometric Growth (2\times):** When the internal size reaches capacity, a new backing array is allocated at double the current size, and existing elements are manually migrated. This keeps individual append operations performing at an average of O(1) time.
* **Asymmetric Shrinking (0.5\times at 1/4 load):** To prevent severe memory bloating, the array down-sizes its backing storage when utilization drops to 25% or less. Shrinking at 1/4 capacity instead of 1/2 is a deliberate architectural pattern designed to prevent **thrashing** (constant allocation/deallocation when an array hovers right on a capacity boundary).

---

## Deep-Dive Lessons Learned & Memory Mechanics

### 1. The Zero-Value Reference Leak Trap
When removing an element from a contiguous array, we shift all subsequent elements to the left and decrement the internal size counter. While this makes the element inaccessible via the array's public APIs, a severe memory bug remains under the hood if left unhandled.

The final, now unused index of the underlying slice still holds a strong reference to the old data. If `T` is a large struct or a pointer, the **Go Garbage Collector (GC)** will track that reference as "live" and refuse to free the memory. We must explicitly overwrite the vacated slot with its type-specific zero value:

```go
var zero T
da.data[da.size] = zero // Breaks the reference chain, allowing immediate GC cleanup

```

### 2. Header Duplication vs. Direct Mutation

Go slices are small descriptor headers composed of a pointer to an underlying array, a length, and a capacity.

Because Go passes all function arguments by value, calling a struct method with a value receiver `func (da DynamicArray[T])` forces the runtime to duplicate this header onto the stack. Any structural adjustments made to `size` or `capacity` inside the method occur exclusively on the copied header and are instantly discarded when the function returns. To enforce sticky updates across data-altering boundaries (`Append`, `Delete`, `grow`), using **pointer receivers** `func (da *DynamicArray[T])` to pass the direct memory address is mandatory.

### 3. Edge-Case Engineering: The Zero-Capacity Initializer

A common vulnerability in custom data structures is the "Zero-Capacity Initialization Bug." If a collection is instantiated with an initial capacity of `0`, the underlying array points to an empty segment.

If your `Append` logic attempts to write directly to `da.data[da.size]` immediately after calling `grow()`, it will crash with a runtime panic if `grow()` updates the pointer *after* the assignment. The growth sequence must completely resolve, allocate a base size of `1`, and swap the underlying data pointer *before* any indexing or assignments take place.

---

## Verification & Test Coverage

To guarantee production readiness, table-driven unit tests were implemented to validate:

* **Boundary Resizing:** Verifying capacity correctly doubles precisely when size crosses the current limit.
* **Panic Defenses:** Catching out-of-bounds requests on negative indices or inputs exceeding `size - 1` without crashing the runtime.
* **GC Integrity:** Asserting that index locations past the current `size` boundary are structurally verified as zeroed out post-deletion.
