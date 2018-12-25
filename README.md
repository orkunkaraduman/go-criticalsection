# DISCONTINUED

Please visit: [github.com/orkunkaraduman/go-syncex](https://github.com/orkunkaraduman/go-syncex)

# Go CriticalSection

[![GoDoc](https://godoc.org/github.com/orkunkaraduman/go-criticalsection?status.svg)](https://godoc.org/github.com/orkunkaraduman/go-criticalsection)

The repository provides `criticalsection` package.

A CriticalSection is a kind of lock like mutex. But it doesn't block
first locked goroutine/section again.

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/orkunkaraduman/go-criticalsection"
)

func main() {
	var cs criticalsection.CriticalSection
	var f int
	cs.Lock()
	for i := 0; i < 5; i++ {
		cs.Lock()
		go func() {
			cs.Lock()
			f++
			fmt.Println("goroutine: ", f)
			cs.Unlock()
		}()
		f++
		fmt.Println("forloop: ", f)
		cs.Unlock()
	}
	fmt.Println("mainfunc: ", f)
	cs.Unlock()
	time.Sleep(1 * time.Second)
	fmt.Println("mainfunc: ", f)
}
```
