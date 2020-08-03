package patterns

import (
	"context"
	"math/rand"
	"fmt"
	"time"
)

func Context(){
	var getStream = func(ctx context.Context) <-chan int{
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

	var ctx, _ = context.WithTimeout(context.Background(), time.Second * 1)
	var stream = getStream(ctx)

	defer fmt.Printf("Terminating...\n")
	for v := range stream {
		fmt.Printf("Value: %d\n", v)
	}
}