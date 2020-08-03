package patterns

import (
	"math/rand"
	"fmt"
	"context"
)

func BridgeChannel(){
	var createChannelOfChannels = func() <-chan chan int{
		var stream = make(chan chan int)

		go func(){
			defer close(stream)

			for i:=0; i < 10; i++ {
				var channel = make(chan int, 1)
				channel <- rand.Intn(100)
				stream <- channel 
				close(channel)
			} 
		}()

		return stream
	}

	var bridge = func(ctx context.Context, channelOfChannels <-chan chan int) <-chan int{
		var stream = make(chan int)

		go func(){
			defer close(stream)

			for {
				select {
					case <-ctx.Done():
						return
					case channel, ok := <-channelOfChannels:
						if !ok {
							return
						}

						for v := range channel {
							stream <- v
						}
				}
			}
		}()

		return stream
	}

	var channelOfChannels = createChannelOfChannels()
	var ctx = context.Background()
	var stream = bridge(ctx, channelOfChannels)

	defer fmt.Printf("Terminating...\n")
	for v := range stream {
		fmt.Printf("Value: %d\n", v)
	}
}