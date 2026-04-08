package main

/**
 * Strategy 1: Shared State with Mutexes
 * Intuition: When multiple goroutines must access the same variable (like a slice),
 * we use a sync.Mutex to "lock" the resource.
 *
 * Performance Tip: Minimize Critical Sections
 * Locking a huge block of code that includes time-consuming logic (like an API call or Sleep)
 * is bad. It forces other goroutines to wait unnecessarily.
 * Instead, perform the "expensive" work first, and ONLY lock during the exact line
 * where you update the shared variable (the "Critical Section").
 */

import (
	"fmt"
	"sync"
)

// confinement pattern
// uses mutex
// thread safe coding. when dealing with multiple concurrent access
var mu sync.Mutex

// this is a race condition code.
// func main() {
// 	start := time.Now()
// 	var arr []int // share resource which everyone wants to write to.
// 	var wg sync.WaitGroup
// 	for i := 1; i <= 5; i++ {
// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup, val int, x *[]int) {
// 			// so for performance optimization
// 			// you onnly lock the section where you are updating
// 			// critical section, not the whole logic, that logic
// 			// can take time, which will block other ideal runnable
// 			// gorotuines from taking the lock and executing.

// 			// mean while what you can do is just lock the critical section
// 			// while the time taking thing can happen concurrently.

// 			// mu.Lock() on critical section has much higer performance than
// 			// mu.Lock() on the entire logic of code.

// 			defer wg.Done()
// 			// mu.Lock()
// 			// time.Sleep(3 * time.Second)
// 			// mu.Lock()
// 			val = val * 2
// 			*x = append(*x, val)
// 			// mu.Unlock()
// 			// mu.Unlock()
// 		}(&wg, i, &arr)
// 	}
// 	wg.Wait()
// 	fmt.Println(arr)
// 	fmt.Println(time.Since(start))
// }

/**
 * Strategy 2: Confinement
 * Intuition: "Don't communicate by sharing memory; share memory by communicating."
 * (Or in this case, partition the memory so you don't have to share it at all).
 *
 * By pre-allocating the slice and passing a pointer to a UNIQUE index to each goroutine,
 * we ensure that no two goroutines ever touch the same memory address.
 *
 * Benefits:
 * 1. Zero Lock Contention: Since there are no locks, there is no waiting.
 * 2. Performance: This is significantly faster than using a Mutex for every operation.
 * 3. Simplicity: You prevent race conditions by design, not by careful locking.
 */
func main() {
	// Pre-allocation: We know exactly how much space we need.
	var arr []int = make([]int, 5)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		// We pass the address of a SPECIFIC index (&arr[i]) to each goroutine.
		// This "confines" the goroutine's scope to just that element.
		go func(wg *sync.WaitGroup, val *int, i int) {
			defer wg.Done()
			
			// Thinking Section: Imagine doing complex math here. 
			// Because you are working on your own private pointer, 
			// you don't need a lock!
			*val = i
		}(&wg, &arr[i], i+1)
	}
	wg.Wait()
	for _, val := range arr {
		fmt.Println(val)
	}
}
