package patterns

import (
	"fmt"
	"math/rand"
	"time"
)

func OrDone(){
	var getStream = func(done <-chan bool) <-chan interface{}{
		var stream = make(chan interface{})
		var timer = time.After(10 * time.Second)

		go func(){
			defer close(stream)

			for {
				select {
				case stream <- rand.Intn(100):
				case <-done:
					return
				case <- timer:
					return
				}
			}
		}()

		return stream
	}

	var orDone = func(done <-chan bool, stream <-chan interface{}) <-chan interface{} {
		var internal = make(chan interface{})

		go func(){
			defer close(internal)

			for {
				select {
				case <-done:
					return
				case v, ok := <-stream:
					if ok {
						internal <- v
					}
				}
			}
		}()

		return internal
	}

	var done = make(chan bool, 2)

	defer fmt.Printf("Terminating...\n")
	defer close(done)

	go func(){
		time.Sleep(2 * time.Second)
		done <- true
		done <- true
	}()

	for v := range orDone(done, getStream(done)) {
		fmt.Printf("Value: %v\n", v)
	}

}