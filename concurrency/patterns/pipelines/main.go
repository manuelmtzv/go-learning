package main

import "fmt"

func generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int, len(integers))
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()
	return multipliedStream
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i + additive:
			}
		}
	}()
	return addedStream
}

// func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
// 	valueStream := make(chan interface{})
// 	go func() {
// 		defer close(valueStream)
// 		for {
// 			for _, v := range values {
// 				select {
// 				case <-done:
// 					return
// 				case valueStream <- v:
// 				}
// 			}
// 		}
// 	}()
// 	return valueStream
// }

// func take(done <-chan interface{}, valueStream <-chan interface{}, limit int) <-chan interface{} {
// 	takeStream := make(chan interface{})
// 	go func() {
// 		defer close(takeStream)
// 		for i := 0; i < limit; i++ {
// 			select {
// 			case <-done:
// 				return
// 			case takeStream <- <-valueStream:
// 			}
// 		}
// 	}()
// 	return takeStream
// }

// func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
// 	valueStream := make(chan interface{})
// 	go func() {
// 		defer close(valueStream)
// 		select {
// 		case <-done:
//			 return
// 		case valueStream <- fn():
// 		}
// 	}()
// 	return valueStream
// }

func main() {
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}