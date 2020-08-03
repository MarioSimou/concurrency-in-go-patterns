package patterns

import (
	"context"
	"runtime"
	"time"
	"fmt"
)

func Pipelines(){
	var pipelineStage = func(ctx context.Context, stream <-chan interface{}) <-chan interface{} {
		var newStream = make(chan interface{})

		go func(){		
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <- stream:
					if !ok {
						return
					}

					time.Sleep(1 * time.Second) // simulate work
					newStream <- v
				}
			}
		}()

		return newStream
	}
	var fanout = func(ctx context.Context, channels ...<-chan interface{}) chan interface{}{
		var stream = make(chan interface{})

		for _, channel := range channels {
			go func(ctx context.Context, channel <-chan interface{}){
				for {
					select {
					case <-ctx.Done():
						return
					case v, ok := <- channel:
						if !ok {
							return
						}
	
						stream <- v
					}
				}
			}(ctx, channel)
		}

		return stream
	}

	var ctx, _ = context.WithTimeout(context.Background(), time.Second * 5)
	var stream = repeat(ctx,1,2,3,4,5)

	var nCPUs = runtime.NumCPU()
	var channels = make([]<-chan interface{}, nCPUs)

	// fan-in = gets a stream and runs it across multiple goroutines
	for i:=0; i < nCPUs; i++ {
		channels[i] = pipelineStage(ctx, stream)
	}

	// fan-out = consumes the data coming from the goroutines and multiplexes it to a single channel
	var fanOutStream = fanout(ctx, channels...)

	defer fmt.Printf("Terminating...\n")
	defer time.Sleep(1 * time.Second)
	defer close(fanOutStream)

	for {
		select {
		case v := <-fanOutStream:
			fmt.Printf("Value: %v\n", v)
		case <-ctx.Done():
			return
		}
	}
}