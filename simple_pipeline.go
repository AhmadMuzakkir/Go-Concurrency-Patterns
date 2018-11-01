package main

import (
	"fmt"
	"sync"
)

// Demonstrate a simple pipeline construction with explicit cancellation.
// From https://blog.golang.org/pipelines

// The pipeline will generate a sequence of integers and then square root each integers.
// We distribute the square root work across two goroutines.

func main() {
	done := make(chan struct{})
	// We can stop the pipeline at any time by closing done.
	defer close(done)

    // Generate a sequence of integers.
	in := generate(done, 1, 10)

	// Distribute the sqroot work across two goroutines that both read from in.
	c1 := sqroot(done, in)
	c2 := sqroot(done, in)

    // Merge the goroutines into one
	out := merge(done, c1, c2)
	for i := range out {
		fmt.Println(i)
	}
}

func generate(done <-chan struct{}, from, to int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := from; i <= to; i++ {
			select {
			case out <- i:
			case <-done:
				return
			}
		}
	}()
	return out
}

func sqroot(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
