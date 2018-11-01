package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)
// Demonstrate how to convert a channel into a blocking function.

// Let's say you have a stream of data in a channel but you don't want to expose the channel across package.
// Instead, you want it to be a blocking function. e.g. data, err := get()

// The function get() should block until there's a new data or an error.
// Error is considered terminal. If there's already an error, get() will immediately return the error.

func main() {
	done := make(chan struct{})
	defer close(done)

	in := generate(done, 1, 5)

	f := newFunctioner(in, done)

	fmt.Println("FIRST")
	for {
		i, err := f.get()
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(i)
		}
	}

	// This will immediately return an error.
	fmt.Println("SECOND")
	for {
		i, err := f.get()
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(i)
		}
	}
}

func generate(done <-chan struct{}, from, to int) <-chan int {
	out := make(chan int, 4)
	go func() {
		defer close(out)
		for i := from; i <= to; i++ {
			select {
			case out <- i:
				time.Sleep(400 * time.Millisecond)
			case <-done:
				return
			}
		}
	}()
	return out
}

type functioner struct {
	in     <-chan int
	signal chan struct{}

	mu  sync.Mutex
	err error
	val int
}

func newFunctioner(in <-chan int, done chan struct{}) *functioner {
	f := &functioner{in: in}
	f.signal = make(chan struct{})

	go func() {
		for n := range f.in {
			f.mu.Lock()
			f.err = nil
			f.val = n
			f.mu.Unlock()

			select {
			case f.signal <- struct{}{}:
			case <-done:
				break
			}
		}

		f.mu.Lock()
		f.err = errors.New("stopped")
		f.val = 0
		f.mu.Unlock()
		close(f.signal)
	}()

	return f
}

// get() will block until there's a new data or an error.
func (f *functioner) get() (int, error) {
	// Check error first. If there's error, immediately return it.
	f.mu.Lock()
	if f.err != nil {
		defer f.mu.Unlock()
		return 0, f.err
	}
	f.mu.Unlock()

	// Block until there's new data or error
	<-f.signal

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.err != nil {
		return 0, f.err
	} else {
		return f.val, nil
	}
}
