package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	for i := 0; i < 3; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				fmt.Println("i", i, "j", j)
			}
		}(i)
	}
	// если посмотреть на вывод, то можно увидеть, как одна из горутин надолго забирает управление (смотрим по i).
	time.Sleep(time.Second)

	fmt.Println("--------------------------------------")
	// runtime.GOMAXPROCS(1) ограничиваем количество ядер до 1
	for i := 0; i < 3; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				fmt.Println("i", i, "j", j)
				runtime.Gosched() // переключаемся на другую горутину
			}
		}(i)
	}
	time.Sleep(time.Second)
}
