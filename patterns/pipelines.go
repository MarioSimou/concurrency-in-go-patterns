package patterns

import (
	"context"
	"fmt"
	"time"
)

func Pipelines(){
	// var repeat = func(ctx context.Context, values ...interface{}) <-chan interface{} {
	// 	var stream = make(chan interface{})

	// 	go func(){
	// 		defer close(stream)

	// 		for {
	// 			for _, v := range values {
	// 				select {
	// 				case <- ctx.Done():
	// 					return
	// 				case stream <- v:
	// 				}
	// 			}
	// 		}
	// 	}()

	// 	return stream
	// } 

	var take = func(ctx context.Context, input <-chan interface{}, n int) <- chan interface{} {
		var stream = make(chan interface{})

		go func(){
			defer close(stream)

			for i :=0; i < n; i++ {
				select {
				case <- ctx.Done():
					return
				case stream <- <-input: 
				} 
			}
		}()

		return stream
	}

	var repeatFn = func(ctx context.Context, fn func() interface{}, n int) <-chan interface{} {
		var stream = make(chan interface{})

		go func(){
			defer close(stream)

			for i:=0; i < n; i++ {
				select {
				case <- ctx.Done():
					return
				case stream <- fn():
				}
			}
		}()

		return stream
	}

	var ctx, cancelFn = context.WithCancel(context.Background())
	// var stream = repeat(ctx, 1, 2, 3, 4, 5)
	var fn = func() interface{} { return 1 }
	var stream = repeatFn(ctx, fn, 10)
	var takeStream = take(ctx, stream, 5)
	var timer = time.After(time.Second * 10)

	defer fmt.Printf("Terminating...\n")
	defer cancelFn()

	for {
		select {
		case <-timer:
			return
		case v, ok := <- takeStream:
			if !ok {
				return
			}

			fmt.Printf("Value: %v\n", v)
		}
	}
}