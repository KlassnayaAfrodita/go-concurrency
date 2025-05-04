package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ { // будет выводиться в разнобой, потому что вызывается 10 горутин. когда управление доходит до конкретной горутины, итератор уже увеличился
		go func() {
			fmt.Println(i)
		}()
	}
	time.Sleep(time.Second / 2)

	fmt.Println("-----------------------------")

	for i := 0; i < 10; i++ { // правильно делать так
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	time.Sleep(time.Second / 2)

}
