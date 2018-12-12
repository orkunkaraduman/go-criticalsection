// +build ignore

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
	cs.Unlock()
	time.Sleep(1 * time.Second)
	fmt.Println("mainfunc: ", f)
}
