package main

import (
	"fmt"
	"time"
)

/*

   go routine leak := prevent go routine running infintely.

*/

// on production servers which run for a very long time.
// if you leave a goroutine which you dont want to run indefintely
// running indefinitely in the background, then it is consuming resources of your computer in background.
// that is the example of a goroutine leak.

// to prevent this we have done pattern.
// allow parent goroutine to cancel its children.

// func main() {

// 	go func() {
// 		for {
// 			select {
// 			default:
// 				fmt.Println("infinitely running go routine because this default condition alwasys succeeds which is the only thing that matters to select statement")
// 			}
// 		}
// 	}()

// 	time.Sleep(5 * time.Second)
// }

func main() {
	done := make(chan bool)
	go func(done <-chan bool) {
		// receive only channel, means here you could only take out values from it
		// not write into the done channel here.
		for {
			select {
			// if we did not use select then how would we able to check if <-done for this
			// because of this for select pattern becomes very important.
			case <-done:
				fmt.Println("coming out of infinitely running goroutine in background")
				return
			default:
				fmt.Println("running goroutine indefinitely in background")
			}
		}
	}(done)
	// child go routine has already started
	// parent goroutine sleeps for 2 seconds
	time.Sleep(2 * time.Second)
	fmt.Println("time to stop it baby")
	// this will give in the signal to <- done in the select and it will stop.
	defer close(done)
}
