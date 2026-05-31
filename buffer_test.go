package lastnbuffer

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewRoundsUpPowerOfTwo(t *testing.T) {
	b := New[int](10)
	if got, want := b.Cap(), 16; got != want {
		t.Fatalf("Cap() = %d, want %d", got, want)
	}
}

func TestNewZeroBecomesOne(t *testing.T) {
	b := New[int](0)
	if got, want := b.Cap(), 1; got != want {
		t.Fatalf("Cap() = %d, want %d", got, want)
	}
}

func TestSnapshotOrderAndOverwrite(t *testing.T) {
	b := New[int](4)
	for i := 1; i <= 6; i++ {
		b.Push(i)
	}

	got := b.Snapshot()
	want := []int{3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Snapshot() = %v, want %v", got, want)
	}
}

func TestLen(t *testing.T) {
	b := New[int](4)
	if got, want := b.Len(), 0; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	for i := 0; i < 10; i++ {
		b.Push(i)
	}

	if got, want := b.Len(), 4; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}
}

func TestRangeOrderAndEarlyStop(t *testing.T) {
	b := New[int](8)
	for i := 1; i <= 5; i++ {
		b.Push(i)
	}

	var got []int
	b.Range(func(v int) bool {
		got = append(got, v)
		return v < 3
	})

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Range result = %v, want %v", got, want)
	}
}

func TestGetLastN(t *testing.T) {
	b := New[int](8)
	for i := 1; i <= 6; i++ {
		b.Push(i)
	}

	if got := b.GetLastN(0); got != nil {
		t.Fatalf("GetLastN(0) = %v, want nil", got)
	}

	got := b.GetLastN(3)
	want := []int{4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetLastN(3) = %v, want %v", got, want)
	}

	got = b.GetLastN(20)
	want = []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetLastN(20) = %v, want %v", got, want)
	}
}

func TestIter(t *testing.T) {
	b := New[int](4)
	for i := 1; i <= 6; i++ {
		b.Push(i)
	}

	var got []int
	for v := range b.Iter() {
		got = append(got, v)
	}

	want := []int{3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Iter() = %v, want %v", got, want)
	}
}

func TestGetLastNInto(t *testing.T) {
	b := New[int](8)
	for i := 1; i <= 6; i++ {
		b.Push(i)
	}

	dst := make([]int, 0, 10)
	got := b.GetLastNInto(dst)
	want := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetLastNInto(dst) = %v, want %v", got, want)
	}

	if cap(got) != cap(dst) {
		t.Fatalf("expected capacity reuse: got cap %d, want %d", cap(got), cap(dst))
	}

	b2 := New[int](16)
	for i := 1; i <= 12; i++ {
		b2.Push(i)
	}

	small := make([]int, 0, 4)
	got2 := b2.GetLastNInto(small)
	want2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	if !reflect.DeepEqual(got2, want2) {
		t.Fatalf("GetLastNInto(small) = %v, want %v", got2, want2)
	}

	if cap(got2) < len(want2) {
		t.Fatalf("expected reallocation with enough capacity, got cap %d", cap(got2))
	}
}

func TestConcurrentPushAndRead(t *testing.T) {
	b := New[int](64)

	const writers = 8
	const perWriter = 2000

	var wg sync.WaitGroup
	wg.Add(writers)

	for w := 0; w < writers; w++ {
		go func(base int) {
			defer wg.Done()
			for i := 0; i < perWriter; i++ {
				b.Push(base*perWriter + i)
				_ = b.Len()
				_ = b.Snapshot()
				b.Range(func(int) bool { return true })
			}
		}(w)
	}

	wg.Wait()

	if got, want := b.Len(), 64; got != want {
		t.Fatalf("Len() after concurrent writes = %d, want %d", got, want)
	}

	if len(b.Snapshot()) != b.Len() {
		t.Fatalf("Snapshot length mismatch")
	}
}
