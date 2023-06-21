package lastnbuffer

import (
	"errors"
	"sync/atomic"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isPowerOfTwo(x uint64) bool {
	return x&(x-1) == 0
}

var (
	ErrInvalidBufferSize = errors.New("buffer must be of size 2^n")
)

type LastNBuffer[T any] struct {
	data          []T
	pos           uint64
	bitWiseLength uint64
}

func NewLastNBuffer[T any](size int) (*LastNBuffer[T], error) {
	if size == 0 || !isPowerOfTwo(uint64(size)) {
		return nil, ErrInvalidBufferSize
	}
	return &LastNBuffer[T]{
		data:          make([]T, size),
		pos:           0,
		bitWiseLength: uint64(size - 1),
	}, nil
}

func (s *LastNBuffer[T]) Reset() {
	var v T
	s.data[0] = v
	atomic.StoreUint64(&s.pos, 0)
}

func (s *LastNBuffer[T]) Append(v T) (old T) {
	index := atomic.AddUint64(&s.pos, 1) - 1
	index &= s.bitWiseLength
	_ = s.data[index]
	old = s.data[index]
	s.data[index] = v
	return old
}

func (s *LastNBuffer[T]) ForEach(fn func(e T) bool) {
	pos := atomic.LoadUint64(&s.pos)

	if int(pos) > len(s.data) {
		pos &= s.bitWiseLength

		data := s.data[pos:]
		for i := 0; i < len(data); i++ {
			if !fn(data[i]) {
				return
			}
		}
	}

	data := s.data[:pos]
	for i := 0; i < len(data); i++ {
		if !fn(data[i]) {
			return
		}
	}
}

func (s *LastNBuffer[T]) GetLastN(n int) []T {
	if n <= 0 {
		n = len(s.data)
	} else {
		n = min(n, len(s.data))
	}

	pos := atomic.LoadUint64(&s.pos)
	ret := make([]T, min(n, int(pos)))
	n = 0
	if int(pos) >= len(s.data) {
		pos &= s.bitWiseLength
		n = copy(ret, s.data[pos:])
	}
	if n < len(ret) {
		copy(ret[n:], s.data[:pos])
	}
	return ret
}

func (s *LastNBuffer[T]) GetAll() []T {
	return s.GetLastN(-1)
}
