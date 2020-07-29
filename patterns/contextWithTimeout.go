package patterns

import (
	"context"
	"fmt"
	"time"
)

func ContextWithTimeouts(){
	type MyError struct {
		Err error
		Time time.Time
	}

	var handleRequest = func(ctx context.Context, errStream chan<-MyError){
		
		time.Sleep(time.Second * 5)

		if e := ctx.Err(); e != nil {
			var deadline, _ = ctx.Deadline()
			errStream <- MyError{Err: ctx.Err(), Time: deadline}
			close(errStream)
		}
	}

	defer fmt.Printf("Terminating...\n")

	var ctx, _ = context.WithTimeout(context.Background(), time.Second * 4)
	var errStream = make(chan MyError)

	go handleRequest(ctx, errStream)

	for e := range errStream {
		fmt.Printf("Response: %v\n", e)
	} 
}