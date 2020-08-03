package patterns

import (
	"time"
	"math/rand"
	"fmt"
)

func getStream() <-chan int {
	var stream = make(chan int)

	go func(){
		var timer = time.After(time.Second * 5)
		defer close(stream)

		for {
			select {
			case <-timer:
				return
			case stream <- rand.Intn(100):
			}
		}
	}()

	return stream
}


func LexicalConfinement(){
	var stream = getStream()

	defer fmt.Printf("Terminating...\n")
	for value := range stream {
		fmt.Printf("Value: %d\n", value)
	}
}