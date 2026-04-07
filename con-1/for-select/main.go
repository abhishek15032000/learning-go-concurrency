package main

import "fmt"

func main() {
	// buffered channel, we wont need a receiver until we hit the capacity
	// for the buffered channel.
	// for making communication between channels asynchronous := we used buffered channels
	// for making communication between channel synchronous we used unbuffered channel
	// make(chan string) with no capacity specified, will make it unbuffered channel
	charChannel := make(chan rune, 3)
	arr := []rune{'a', 'b', 'c'}
	for _, val := range arr {
		select {
		// why use this, because it would help in making something called generator
		// which generates data, like a producer
		case charChannel <- val:
		}
	}
	close(charChannel)
	for val := range charChannel {
		fmt.Println(string(val))
	}

	// go func() {
	// 	for {
	// 		select {
	// 		default:
	// 			fmt.Println("infinitely running go routine because this default condition alwasys succeeds which is the only thing that matters to select statement")
	// 		}
	// 	}
	// }()
}
