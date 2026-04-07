package main

import (
	"fmt"
)

func someFunc(num int) {
	fmt.Println(num)
}

// my thought process, how i am visualizing goroutines.

// right now the code runs synchronously main cannot exit without someFunc() returning.

// func main() {
// 	someFunc(1)
// 	fmt.Println("hello")
// }

// but if i makethe someFunc a go routine - which means i forking a new prcoess from the main process
// then this wont print 1, because the process we forked off did not join the main process before
// the main process exited.

// func main() {
// 	go someFunc(1)
// 	fmt.Println("Hello")
// }

// now what we can do is like block the main process or make it wait so that process someFunc() would
// join it at some point of time.
// but the problem here is the time we expect that by this much time it could join, it maybe
// too large, or to less so our output becomes unpredicatable.

// func main() {
// 	go someFunc(1)
// 	time.Sleep(1 * time.Microsecond) // make it time.Second we get an output.
// 	fmt.Println("Hello")
// }

/*
 	these process which are forked off are asynchronous
	they dont wait for other process to finish off, to start their own.
	as you can see the output,the output is unpredictable since
	we know every process will join main process but we dont know the order
	in which they will be executed. because go scheduler schedules them.
*/

// func main() {

// 	go someFunc(1)
// 	go someFunc(2)
// 	go someFunc(3)
// 	time.Sleep(2 * time.Second)
// 	fmt.Println("Hello")
// }

// channels := if two goroutines want to share information between each other they need channels
// all go routines run independently, they dont require knowledge about other go routines running.
// channels are first-in-first-out queues.

// implementing a join point.
// this makes the anyonymous process and main process becomes synced because
// unbuffered channels are blocking.

// func main() {

// 	mychannel := make(chan string)
// 	// anonymous function
// 	go func() {
// 		mychannel <- "data"
// 	}()
// 	// created joining point.
// 	msg := <-mychannel // this line of code is blocking, main process is waiting either
// 	// to receive in the channel or the channel gets close.
// 	fmt.Println(msg)
// }

// select statement := it helps a goroutine wait on multiple communication operations.

// select is a blocking piece of code.
// which means it will block the main process until any one of the cases gets executed.
// if there are messages in multiple channels at the same time ready to be taken
// it chooses one randomly, and then it unblocks the main process.

// the process happens something like this <- data happens it is waiting for the
// receiver to be ready, as soon as the receiver is ready transfer happens and then
// unblocks the process.

func main() {
	chan1 := make(chan string)
	chan2 := make(chan string)
	go func() {
		chan1 <- "data"
		close(chan1) // always the sender should close the channel.
	}()
	go func() {
		chan2 <- "hello"
		close(chan2)
	}()

	select {
	case val, ok := <-chan1:
		if !ok {
			fmt.Println("channel 1 closed")
		}
		fmt.Println(val)
	case val, ok := <-chan2:
		if !ok {
			fmt.Println("channel 2 closed")
		}
		fmt.Println(val)
	}
}
