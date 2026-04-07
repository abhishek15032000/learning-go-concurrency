package main

import "fmt"

// stage1 channel - seperation of concern
func sliceToChannel(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		//asynchronous, spawn goroutine but slicetochannel wont wait for this to complete its entire execution

		defer close(out)
		for _, val := range nums {
			out <- val
		}
	}()
	return out
}

// stage2 channel - seperation of concern
func squareit(val <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		// looping until the channel val closes.
		// for a read only channel we have to range loop
		for x := range val {
			out <- (x * x)
		}
	}()
	return out
}

func main() {
	// input

	nums := []int{2, 3, 4, 7, 15}
	// stage1 := add each value in this array to the channel
	datachannel := sliceToChannel(nums)
	// stage2 := pass values in data channel to sq channel which will square
	// the nums.
	sqchannel := squareit(datachannel)
	// stage3 := print the squared up nos in a string
	for val := range sqchannel {
		fmt.Println(val)
	}
}
