package lastnbuffer

import (
	"iter"
	"sync"
)

type Buffer[T any] struct {
	mu    sync.RWMutex
	data  []T
	mask  uint64
	write uint64
}

// New creates a thread-safe ring buffer that stores the last N elements.
//
// The requested capacity n is rounded up to the next power of two.
// If n is 0, capacity becomes 1.
func New[T any](n uint64) *Buffer[T] {
	capPow2 := nextPow2(n)
	return &Buffer[T]{
		data: make([]T, capPow2),
		mask: capPow2 - 1,
	}
}

// Push appends one element to the buffer.
//
// When the buffer is full, the oldest element is overwritten.
func (b *Buffer[T]) Push(v T) {
	b.mu.Lock()
	b.data[b.write&b.mask] = v
	b.write++
	b.mu.Unlock()
}

// Len returns how many elements are currently available.
//
// The returned value is in [0, Cap()].
func (b *Buffer[T]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.lenLocked()
}

// Cap returns the effective capacity of the buffer.
//
// Capacity is always a power of two.
func (b *Buffer[T]) Cap() int {
	return len(b.data)
}

// Range iterates elements in chronological order (oldest to newest).
//
// Iteration stops early when fn returns false.
func (b *Buffer[T]) Range(fn func(T) bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	n := b.lenLocked()
	start := b.write - uint64(n)
	for i := 0; i < n; i++ {
		if !fn(b.data[(start+uint64(i))&b.mask]) {
			return
		}
	}
}

// GetLastN returns the last n elements in chronological order.
//
// If n <= 0, it returns nil.
// If n is larger than the number of available elements, it returns all available elements.
func (b *Buffer[T]) GetLastN(n int) []T {
	if n <= 0 {
		return nil
	}
	return b.GetLastNIntoN(nil, n)
}

// GetLastNInto returns all available elements in chronological order using dst as storage.
//
// It reuses dst when capacity is enough; otherwise it allocates a new slice.
// This is equivalent to GetLastNIntoN(dst, b.Len()).
func (b *Buffer[T]) GetLastNInto(dst []T) []T {
	return b.GetLastNIntoN(dst, b.Len())
}

// GetLastNIntoN returns the last n elements in chronological order using dst as storage.
//
// If n <= 0, it returns dst[:0].
// If n is larger than the number of available elements, it returns all available elements.
// It reuses dst when capacity is enough; otherwise it allocates a new slice.
func (b *Buffer[T]) GetLastNIntoN(dst []T, n int) []T {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if n <= 0 {
		return dst[:0]
	}

	available := b.lenLocked()
	if n > available {
		n = available
	}

	if n == 0 {
		return dst[:0]
	}

	if cap(dst) < n {
		dst = make([]T, n)
	} else {
		dst = dst[:n]
	}

	start := b.write - uint64(n)
	for i := 0; i < n; i++ {
		dst[i] = b.data[(start+uint64(i))&b.mask]
	}
	return dst
}

// Iter returns an iterator over available elements in chronological order
// (oldest to newest).
func (b *Buffer[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		b.mu.RLock()
		defer b.mu.RUnlock()

		n := b.lenLocked()
		start := b.write - uint64(n)
		for i := 0; i < n; i++ {
			if !yield(b.data[(start+uint64(i))&b.mask]) {
				return
			}
		}
	}
}

func (b *Buffer[T]) lenLocked() int {
	if b.write >= uint64(len(b.data)) {
		return len(b.data)
	}
	return int(b.write)
}

func nextPow2(n uint64) uint64 {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}
