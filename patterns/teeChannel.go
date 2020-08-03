package patterns

import (
	"context"
	"math/rand"
	"fmt"
	"time"
	
)

func TeeChannel(){
	var getStream = func(ctx context.Context) <-chan int {
		var stream = make(chan int)

		go func(){
			defer close(stream)

			for {
				select {
				case <-ctx.Done():
					return
				case stream <- rand.Intn(100):
				}
			}
		}()

		return stream
	}

	var teeChannel = func(ctx context.Context, stream <-chan int, n int) []chan int {
		var channels = make([]chan int, n)

		for i:=0; i < n; i++ {
			channels[i] = make(chan int,100)
		}

		var closeChannels = func(channels []chan int){
			for _, channel := range channels {
				close(channel)
			}
		}

		go func(){
			defer closeChannels(channels)

			for {
				select {
				case <- ctx.Done():
					return
				case v, ok := <- stream:
					if !ok {
						return
					}

					for _, channel := range channels {
						channel <- v
					}
				}
			}
		}()

		return channels
	}

	var ctx, cancelFn = context.WithCancel(context.Background())
	var stream = getStream(ctx)
	var channels = teeChannel(ctx, stream, 2)
	var timer = time.After(1 * time.Second)

	defer fmt.Printf("Terminating...\n")
	defer time.Sleep(1 * time.Second)
	defer cancelFn()

	for {
		select {
		case v := <-channels[0]:
			fmt.Printf("(%d) Value: %d\n", 0, v)
		case v := <-channels[1]:
			fmt.Printf("(%d) Value: %d\n", 1, v)
		case <- timer:
			return
		}
	}
}