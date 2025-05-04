package main

import "fmt"

func main() {
	c := make(chan int)
	// generate(c)
	//! работать не будет, тк мы переходим в generate, пишем в канал, но дальше двинуться не можем,
	//! тк канал небуфферизированный и его никто не слушает (потому что мы зависли на generate)
	go generate(c)

	for x := range c {
		fmt.Println(x)
	}
	//? альтернатива
	c2 := make(chan int, 5)

	generate(c2)

	for x := range c2 {
		fmt.Println(x)
	}
}

func generate(ch chan<- int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)

}
