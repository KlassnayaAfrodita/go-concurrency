//! pipeline безумно эффективен в комплекте с fan in - fan out
//! например, когда у нас есть несколько этапов, некоторые из которых выполняются быстро, а другие медленно
//! разбиваем медленную задачу с помощью fan in - fan out на несколько воркеров и делаем ее асинхронно

// 		   fanOut - fanIn
// 			      +---+
// +-----+   +--->|   |--->+   +-----+
// |     |   |    +---+    |   |     |
// |     |-->|--->|   |--->+-->|     |
// |     |   |    +---+    |   |     |
// +-----+   +--->|   |--->+   +-----+
// 			      +---+

//? why use - обрабатывать данные последовательно удобнее, чем накопление данных

//* when use - когда можно разбить процесс на этапы, и каждый этап выполнять асинхронно

//! solution - один этап - одна горутина, канал - средство коммуникации между этапами

package main

import (
	"fmt"
	"sync"
	"time"
)

func generateWork(work []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, data := range work {
			out <- data
		}
	}()
	return out
}

func filter(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for data := range in {
			if data%2 == 0 {
				out <- data
			}
		}
	}()
	return out
}

func half(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for data := range in {
			out <- data / 2
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for data := range in {
			out <- data * data
		}
	}()
	return out
}

func longWork(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for data := range in {
			time.Sleep(500 * time.Millisecond)
			out <- data / 3
		}
	}()
	return out
}

func fanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	wg.Add(len(inputs))
	for _, ch := range inputs {
		go func(ch <-chan int) {
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
			out <- data
		}
	}()
	return out
}

func main() {
	work := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	in := generateWork(work) // начало конвейера

	//! выполняем долгую работу
	out1 := fanOut(in)
	out2 := fanOut(in)
	out3 := fanOut(in)
	out4 := fanOut(in)
	out5 := fanOut(in)
	out1 = longWork(out1)
	out2 = longWork(out2)
	out3 = longWork(out3)
	out4 = longWork(out4)
	out5 = longWork(out5)
	out := fanIn(out1, out2, out3, out4, out5)

	out = filter(out)
	out = square(out)

	for data := range out {
		fmt.Println(data)
	}

}
