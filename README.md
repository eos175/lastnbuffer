# lastnbuffer

`lastnbuffer` is a thread-safe generic ring buffer for keeping the last N items.

- Capacity is rounded up to the next power of two.
- Indexing uses a bitmask (`write & (cap-1)`) instead of modulo.
- Reads are returned in chronological order (oldest to newest).

## Install

```bash
go get github.com/eos175/lastnbuffer
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/eos175/lastnbuffer"
)

func main() {
	b := lastnbuffer.New[int](10) // rounded to 16

	for i := 1; i <= 20; i++ {
		b.Push(i)
	}

	fmt.Println(b.Cap())      // 16
	fmt.Println(b.Len())      // 16
	fmt.Println(b.GetLastN(5)) // [16 17 18 19 20]

	b.Range(func(v int) bool {
		fmt.Println(v)
		return true
	})

	for v := range b.Iter() {
		_ = v
	}
}
```

## API

- `New[T any](n uint64) *Buffer[T]`
- `(*Buffer[T]).Push(v T) uint64`
- `(*Buffer[T]).Len() int`
- `(*Buffer[T]).Cap() int`
- `(*Buffer[T]).GetLastN(n int) []T`
- `(*Buffer[T]).GetLastNInto(dst []T) []T`
- `(*Buffer[T]).Range(fn func(T) bool)`
- `(*Buffer[T]).Iter() iter.Seq[T]`
