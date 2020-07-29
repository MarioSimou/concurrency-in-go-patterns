package patterns

import (
	"fmt"
	"math/rand"
	"time"
)
func Tee(){

	var getSteam = func(done <-chan bool) <-chan int{
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
	
	var teeChannel = func(done <-chan bool, stream <-chan int) (<-chan int, <-chan int){
		var c1 = make(chan int)
		var c2 = make(chan int)
	
		go func(){
			defer close(c1)
			defer close(c2)
	
			for {
				select {
				case <-done:
					return
				case v := <- stream:
					c1 <- v
					c2 <- v
				}
			}
		}()
	
		return c1, c2
	}

	var done = make(chan bool, 2) // buffered done channel
	var stream = getSteam(done)
	var c1, c2 = teeChannel(done, stream)
	var timer = time.After(2 * time.Second)


	defer fmt.Printf("Terminating the app\n")
	defer time.Sleep(1 * time.Second)

	
	for {
		select {
		case <-timer:
			done <- true
			done <- true
			return
		case v, ok := <- c1:
			if ok {
				fmt.Printf("Channel 1: %d\n", v)
			}
		case v, ok := <- c2:
			if ok {
				fmt.Printf("Channel 2: %d\n", v)
			}
		}
	}
}
