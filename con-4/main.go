package main

import (
	"fmt"
	"sync"
	"time"
)

// or done pattern, implementing something similar to context package.
var wg sync.WaitGroup

func main() {
	done := make(chan interface{})
	defer close(done) // graceful shutdown of go routines.
	// done <- "string"
	// done <- 24

	// var1 := <-done
	// var2 := <-done
	// fmt.Println(var1, var2)

	cow := make(chan interface{}, 100)
	pigs := make(chan interface{}, 100)

	wg.Add(2)

	// producer1
	go func() {
		for {
			select {
			case <-done:
				close(cow)
				fmt.Println("closing producer channel cows")
				return
			case cow <- "cow":
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				close(pigs)
				fmt.Println("closing producer channel pigs")
				return
			case pigs <- "pigs":
			}
		}
	}()

	x := orDone(done, cow)
	y := orDone(done, pigs)

	// time.Sleep(2 * time.Second)
	// close(done)
	// wg.Wait()
	go consumeCows(x)
	go consumeCows(y)
	time.Sleep(2 * time.Second)
}

func consumeCows(cows <-chan interface{}) {
	defer wg.Done()
	for val := range cows {
		fmt.Println(val)
	}
}

func consumePigs(pigs <-chan interface{}) {
	defer wg.Done()
	for val := range pigs {
		fmt.Println(val)
	}
}

func orDone(done <-chan interface{}, val <-chan interface{}) <-chan interface{} {
	relayStream := make(chan interface{})
	go func() {
		defer close(relayStream)
		for {
			select {
			case <-done:
				return
			case msg, ok := <-val:
				if !ok {
					fmt.Println("closed channel")
					return
				}
				select {
				case relayStream <- msg:
					// in case the receiver is not receiving the value.
					// then if we are listening to done request which would satisfy the select
				case <-done:
					return
				}
			}
		}
	}()
	return relayStream
}
