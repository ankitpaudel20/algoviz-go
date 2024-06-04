package main

// import (
// 	"fmt"
// )

// // ensureNoPanic tries to recover from a panic that might happen during initialization,
// // especially if a non-hashable type causes a panic in internal map operations.
// func ensureNoPanic(test chan int) {

// 	select {
// 	case val := <-test:
// 		fmt.Printf("got %v from the channel", val)
// 	default:
// 		fmt.Print("did not get anything, in default block")
// 	}

// }

// func main() {

// 	test_chan := make(chan int, 1)
// 	go ensureNoPanic(test_chan)
// 	for i := 0; i < 10; i++ {
// 		fmt.Printf("sending %v to the channel\n", i)
// 		test_chan <- i
// 	}

// }
