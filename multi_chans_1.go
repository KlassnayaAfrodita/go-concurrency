//* горутины обмениваются данными через каналы
//* канал синхронизирует исполнение в разных горутинах
//* потокобезосная очередь
//* можно определелять каналы только на чтение или только на запись, запись и чтение - атомарные операции, те либо полностью записалось, либо не записалось вообще
//* могут быть буфферизированные (1) или небуфферизированными (2). (2) - если хочет прочитать, то ждет, когда другая горутина в него напишет
//*

package main

import "fmt"

func generate(in chan<- int) {
	for i := 0; i < 5; i++ {
		in <- i
	}
	close(in)
}

func square(in <-chan int, out chan<- int) { //? пример синтаксиса
	for i := range in {
		out <- i * i
	}
	close(out)
}

func main() {
	chanIn := make(chan int)
	chanOut := make(chan int)

	go generate(chanIn)
	go square(chanIn, chanOut)

	for i := range chanOut {
		fmt.Println(i)
	}
}
