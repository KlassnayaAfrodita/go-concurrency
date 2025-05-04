package main

import (
	"fmt"
	"sync"
)

func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	wg.Add(len(inputs))
	for _, ch := range inputs {
		go func(ch <-chan int) { //! нужно принимать определенный канал, потому что пока мы проходимся по каналу
			//! внутри горутины, делаем другую логику; этот же канал берет другая горутина
			//! мы получаем состояние гонки и ошибку при инкриминтировании. чтобы избежать - делаем так
			for {
				value, ok := <-ch

				if !ok {
					defer wg.Done()
					break
				}

				out <- value
			}
		}(ch)
	}

	go func() {

		wg.Wait()
		close(out)
	}()

	return out
}

func fanOut(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {

		defer close(out)

		for data := range in {
			fmt.Printf("%p, %d\n", out, data)
			out <- data
		}
	}()
	return out
}

func generateWork(work []int) <-chan int {
	out := make(chan int)
	go func() {

		defer close(out)

		for _, i := range work {
			out <- i
		}
	}()
	return out
}

func main() {
	work := []int{1, 2, 3, 4, 5, 6, 7, 8}

	workChannel := generateWork(work)

	out1 := fanOut(workChannel)
	out2 := fanOut(workChannel)
	out3 := fanOut(workChannel)

	out := fanIn(out1, out2, out3)
	for i := range out {
		fmt.Println(i)
	}
}
