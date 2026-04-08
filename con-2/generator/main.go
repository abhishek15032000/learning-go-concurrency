package main

/**
 * Intuition: Concurrency Pipelining, Fan-out, and Fan-in
 *
 * This program demonstrates a highly efficient way to process data in Go using "pipelines".
 * Instead of doing all work in a single loop, we break the task into independent stages
 * that communicate via channels.
 *
 * Parallelism (Fan-out/Fan-in):
 * Some stages (like prime testing) are computationally expensive (CPU-bound). We can
 * "fan-out" this stage by starting multiple goroutines (one for each CPU core) to
 * process numbers in parallel. We then "fan-in" the results back into a single channel.
 *
 * Industry Use Cases:
 * 1. Data ETL Pipelines: Fetching, transforming, and loading large datasets.
 * 2. Image/Video Processing: Applying filters to frames in parallel.
 * 3. Web Scraping: Fetching multiple URLs simultaneously and processing their content.
 * 4. Microservices: Aggregating data from multiple downstream services.
 */

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// repeatFunc is a "Generator". It turns a standard function (fn) into a
// stream of data on a channel. This is the starting point of our pipeline.
func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	// fn is a first class function, we can pass function to other functions as input.
	// Adding a buffer introduces "Buffered Backpressure".
	// This allows the producer to keep working even if the consumer is temporarily slow,
	// effectively "smoothing out" bursty performance. However, if the consumer is
	// consistently slower, the buffer will fill up and the producer will still block.
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				fmt.Println("stop generating data")
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

// take is a "Pipeline Utility" that limits the flow. It reads 'n' items from a
// stream and then stops, preventing the infinite generator from running forever.
func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()
	return taken
}

// primeFinder is a "Pipeline Stage". It filters the input stream, only passing
// through numbers that are prime. This is the computationally expensive part.
func primeFinder(stream <-chan int, done <-chan any) <-chan int {
	isPrime := func(randomInt int) bool {
		if randomInt < 2 {
			return false
		}
		for i := 2; i*i <= randomInt; i++ {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case val := <-stream:
				if isPrime(val) {
					primes <- val
				}
			}
		}
	}()
	return primes
}

func main() {
	start := time.Now()
	done := make(chan interface{})
	defer close(done)

	randomNumberFetcher := func() int {
		return rand.Intn(500_000)
	}

	// Stage 1: Generate an infinite stream of random numbers.
	x := repeatFunc(done, randomNumberFetcher)

	// Stage 2: Fan-out (Parallelize).
	// We check how many CPU cores we have and spin up that many primeFinders.
	// This allows us to utilize all available hardware power.
	cores := runtime.NumCPU()
	arr := make([]<-chan int, cores)
	for i := 0; i < cores; i++ {
		arr[i] = primeFinder(x, done)
	}

	// Stage 3: Fan-in.
	// We merge the 'n' multiple prime channels back into a single result stream.
	y := fanIn(done, arr)

	// Stage 4: Take.
	// We only want the first 10 primes found.
	for val := range take(done, y, 10) {
		fmt.Println(val)
	}
	fmt.Println("Time used := ", time.Since(start))
}

// fanIn is a multiplexer. It joins multiple channels into one.
// It uses a sync.WaitGroup to ensure the output channel is closed only after
// ALL input channels have finished their work.
func fanIn(done <-chan any, arr []<-chan int) <-chan int {
	out_stream := make(chan int)
	var wg sync.WaitGroup
	// starting multiple workers:= worker pool to do this
	for _, val := range arr {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case nums := <-val:
					out_stream <- nums
				}
			}
		}(&wg)
	}
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(out_stream)
	}(&wg)
	return out_stream
}
