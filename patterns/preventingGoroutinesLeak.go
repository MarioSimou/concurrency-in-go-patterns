package patterns

import (
	"math/rand"
	"time"
	"fmt"
)

func PreventingGoroutinesLeak(){
	var getStream = func(done <-chan bool) <-chan int {
		var stream = make(chan int)

		go func(){
			defer close(stream)

			for {
				select {
				case <-done:
					fmt.Printf("Closing done channel...\n")
					return
				case stream <- rand.Intn(100):
				}
			}
		}()

		return stream
	}
	
	var done = make(chan bool,1)
	var stream = getStream(done)
	var timer = time.After(time.Second)

	defer fmt.Printf("Terminating...\n")
	defer close(done)

	for {
		select {
		case v := <-stream:
			fmt.Printf("Value: %d\n", v)
		case <-timer:
			done <- true
			return
		}
	}
}