package patterns

import (
	"math/rand"
	"time"
	"fmt"
)

func OrDoneChannel(){
	var getStream = func(done <-chan bool) <-chan int {
		var stream = make(chan int)

		go func(){
			defer close(stream)

			for {
				select {
				case <-done:
					return
				case stream <- rand.Intn(100):
				}
			}
		}()

		return stream
	}

	var orDone = func(done <-chan bool, stream <-chan int) <-chan int{
		var newStream = make(chan int)

		go func(){
			defer close(newStream)

			for {
				select {
				case <-done:
					return
				case newStream <- <-stream:
				}
			}
		}()

		return newStream
	}

	var done = make(chan bool, 2)
	var stream = getStream(done)

	go func(){
		time.Sleep(1 * time.Second)
		done <- true
		done <- true
	}()

	defer fmt.Printf("Terminating...\n")
	defer time.Sleep(1 * time.Second)
	defer close(done)
	
	for v := range orDone(done, stream) {
		fmt.Printf("Value: %d\n", v)
	}
}