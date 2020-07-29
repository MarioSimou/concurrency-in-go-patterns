package patterns

import (
	"fmt"
	"math/rand"
)

// Bridge pattern gets a channel of channels as an input and return a channel. It orchestrates how values are passed from multiple channels
// to a single channel
func Bridge(){
	var n = 10

	// Generates 10 channels and pass a random number to each one
	var getStream = func() <-chan chan int {
		var stream = make(chan chan int, n)
		
		go func(){
			defer close(stream)

			for i := 0; i < n; i++ {
				var c = make(chan int, 1)
				c <- rand.Intn(100)
				stream <- c
				close(c)	
			}

		}()

		return stream
	}

	var getBridge = func(done <-chan bool, nested <-chan chan int) <-chan int {
		var stream = make(chan int, n)

		go func(){
			defer close(stream)

			for {
				select {
				case <- done:
					return
				case c, ok := <-nested:
					if !ok {
						return
					}

					for v := range c {
						stream <- v
					}

				}
			}
		}()

		return stream
	}


	defer fmt.Printf("Terminating...\n")

	var done = make(chan bool, n)
	var stream = getStream()
	var bridge = getBridge(done, stream)

	defer close(done)

	for v := range bridge {
		fmt.Printf("Value: %d\n", v)
	}
}