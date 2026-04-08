package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// context package.

// helps preventing goroutine leaks

// a context is immutable, it cannot be changed by its children go routines
// it cannot be canceled by children go routines
// the only thing passed to children is the context, not the cancel function.
// because we need to maintain clarity on who is orcestrating this stuff

// provides api for cancelling out aysnchronous goroutines.
func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	generator := func(dataItem string, stream chan interface{}) {
		for {
			select {
			case <-ctx.Done():
				close(stream)
				fmt.Println("closing producer of ", dataItem)
			case stream <- dataItem:
			}
		}
	}
	infiniteApples := make(chan interface{})
	go generator("apple", infiniteApples)
	infiniteOranges := make(chan interface{})
	go generator("oranges", infiniteOranges)
	infinitePeaches := make(chan interface{})
	go generator("peaches", infinitePeaches)

	wg.Add(3)
	go func1(ctx, &wg, infiniteApples)
	func2 := genericFunc
	func3 := genericFunc
	go func2(ctx, &wg, infiniteOranges)
	go func3(ctx, &wg, infinitePeaches)

	time.Sleep(3 * time.Second)
}

func func1(ctx context.Context, wg *sync.WaitGroup, stream chan interface{}) {
	defer wg.Done()
	var wg2 sync.WaitGroup
	doWork := func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-stream:
				if !ok {
					fmt.Println("channel closed")
					return
				}
				fmt.Println(d)
			}
		}
	}
	newContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	for i := 0; i < 3; i++ {
		wg2.Add(1)
		go doWork(newContext, &wg2)
	}
	defer cancel()
	wg2.Wait()
}

func genericFunc(ctx context.Context, wg *sync.WaitGroup, stream <-chan interface{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-stream:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println(d)
		}
	}
}
