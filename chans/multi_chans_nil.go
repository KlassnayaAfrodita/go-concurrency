package main

import "fmt"

func main() {
	var a chan int

	// a <- 1
	// fmt.Println(<-a)

	select {
	case i := <-a:
		fmt.Println(i)
	default:
		fmt.Println("default")
	}
	close(a) //! нельзя закрывать пустой канал, тк мы его объявили, но не инициализировали
}

//! если канал определяем, то его надо инициализировать через make!!!!!!!!!!!!!!
