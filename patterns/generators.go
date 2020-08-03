package patterns

import (
	"context"
	"fmt"
	"time"
)

func repeat(ctx context.Context, args ...interface{}) <-chan interface{} {
	var stream = make(chan interface{})

	go func(){
		defer close(stream)

		for {
			for _, arg := range args {
				select {
				case <-ctx.Done():
					return
				case stream <- arg:
				}
			}
		}
	}()

	return stream
}

func take(ctx context.Context, stream <-chan interface{}, n int) <-chan interface{} {
	var newStream = make(chan interface{})

	go func(){
		defer close(newStream)

		for i:=0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-stream:
				if ok {
					newStream <- v
				}
			}
		}
	}()

	return newStream
}

func Generators() {
	var ctx, _ = context.WithTimeout(context.Background(), time.Millisecond * 1)
	var stream = take(ctx, repeat(ctx, 1,2,3,4,5), 10)

	defer fmt.Printf("Terminating...\n")
	for {
		select {
		case v, ok := <-stream:
			if !ok {
				return
			}
			
			fmt.Printf("Value: %v\n", v)
		case <-ctx.Done():
			return
		}
	}
}