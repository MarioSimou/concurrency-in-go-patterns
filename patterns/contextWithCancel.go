package patterns 

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func ContextWithCancel(){
	var generateInts = func(ctx context.Context) <-chan int{
		var stream = make(chan int)

		go func(){
			defer close(stream)

			for {
				select {
				case <- ctx.Done():
					return
				case stream <- rand.Intn(100):  
				}
			}

		}()

		return stream
	}

	var timer = time.After(3 * time.Second)
	var ctx, cancelFn = context.WithCancel(context.Background())
	var stream = generateInts(ctx)

	defer fmt.Printf("Terminating...\n")

	for {
		select {
		case <-timer:
			cancelFn()
		case v, ok := <- stream:
			if !ok {
				return
			}

			fmt.Printf("Value: %d\n", v)
		}
	}

}