package patterns

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func FanOutFanIn(){
	var findPrime = func(ctx context.Context) <-chan interface {} {
		var n = rand.Intn(1000)
		var timer = time.After(time.Duration(n) * time.Millisecond)
		var stream = make(chan interface{})

		go func(){
			defer close(stream)

			for {
				select {
				case <-timer:
					stream <- n
					return
				case <- ctx.Done():
					return
				}
			}
		}()

		return stream
	}

	var fanIn = func(ctx context.Context, streams ...<-chan interface{}) <- chan interface{} {
		var stream = make(chan interface{})

		go func(){
			defer close(stream)

			for _, s := range streams {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-s:
					if !ok {
						continue
					}
					stream <- v
				}
			}
		}()

		return stream
	}

	var take = func(ctx context.Context, stream <-chan interface{}, n int) <-chan interface{} {
		var newStream = make(chan interface{})

		go func(){
			defer close(newStream)

			for i:=0; i < n; i++ {
				select {
				case <-ctx.Done():
				case v, ok := <-stream:
					if ok {
						newStream <- v
					}
				}
			}
		}()

		return newStream
	}

	var ctx, cancelFn = context.WithCancel(context.Background())
	var nCPUs = runtime.NumCPU()
	var allChannels = make([]<-chan interface{}, nCPUs)

	for i := 0; i < nCPUs; i++ {
		allChannels[i] = findPrime(ctx)
	}

	var fanInStream = fanIn(ctx, allChannels...)

	defer fmt.Printf("Terminating....\n")
	defer cancelFn()

	for v := range take(ctx, fanInStream, 6) {
		fmt.Printf("Value: %v\n", v)
	}
}