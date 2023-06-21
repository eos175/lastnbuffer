# LastNBuffer

The `LastNBuffer` package provides a concurrent circular buffer implementation that stores the last N elements of a given type. It is useful in scenarios where you need to maintain a fixed-size buffer and only keep track of the most recent elements while discarding the oldest ones.


## Usage

To start using the `LastNBuffer` package, follow these steps:

1. Import the package into your Go code:
```go
import "github.com/eos175/lastnbuffer"
```

2. Create a new `LastNBuffer` instance with a specified buffer size:
```go
bufferSize := 16
buffer, err := lastnbuffer.NewLastNBuffer[string](bufferSize)
if err != nil {
    // Handle error
}
```

3. Append elements to the buffer using the `Append` method:
```go
element := "example"
oldElement := buffer.Append(element)
```

4. Retrieve the last N elements from the buffer using the `GetLastN` method:
```go
elements := buffer.GetLastN(5)
```

5. Iterate over the elements in the buffer using the `ForEach` method:
```go
buffer.ForEach(func(x string) bool {
    // Process x
    return true // Continue iteration
})
```

6. Reset the buffer to its initial state using the `Reset` method:
```go
buffer.Reset()
```

7. Get all elements in the buffer using the `GetAll` method:
```go
data := buffer.GetAll()
```

## Example

Here's a simple example demonstrating the usage of the `LastNBuffer` package:

```go
package main

import (
	"fmt"
	"github.com/eos175/lastnbuffer"
)

func main() {
	bufferSize := 8
	buffer, _ := lastnbuffer.NewLastNBuffer[int](bufferSize)


	buffer.Append(1)
	buffer.Append(2)
	buffer.Append(3)
	buffer.Append(4)
	buffer.Append(5)

	elements := buffer.GetLastN(3)
	fmt.Println("Last 3 elements:", elements) // Output: [3 4 5]

	buffer.ForEach(func(element int) bool {
		fmt.Println("Processing element:", element)
		return true
	})

	buffer.Reset()
	elements = buffer.GetAll()
	fmt.Println("All elements:", elements) // Output: []
}
```

## Concurrency

The `LastNBuffer` implementation is concurrent-safe due to the use of the `atomic` package in Go. The methods provided by the `LastNBuffer` type ensure that the buffer operations are atomic and can be safely accessed by multiple goroutines.

## Error Handling

The `NewLastNBuffer` function returns an error (`ErrInvalidBufferSize`) if the provided buffer size is not a power of two or zero. You should handle this error appropriately in your code.

## License

This project is licensed under the [MIT License](LICENSE).
