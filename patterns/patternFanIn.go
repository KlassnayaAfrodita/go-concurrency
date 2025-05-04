//? why use - чтобы объеденить данные из нескольких каналов

//* when use - когда нучто объеденить данные из нескольких каналов

//! solution - проходим по списку каналов и запускаем горутины, в которых сливаем все в один канал

package main

import (
	"fmt"
	"sync"
)

func fanIn(channels ...<-chan int) <-chan int { // принимаем на вход много каналов
	sink := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(channels)) // пополняем вейтгруппу на количество каналов

	for _, channel := range channels { // идем по каналам
		go func(ch <-chan int) { //! нужно принимать определенный канал, потому что пока мы проходимся по каналу
			//! внутри горутины, делаем другую логику; этот же канал берет другая горутина
			//! мы получаем состояние гонки и ошибку при инкриминтировании. чтобы избежать - делаем так
			// запускаем отдельно горутины, чтобы писать в общий канал

			for { // проходимся по каналу, пока он не закончятся значения в канале
				value, ok := <-ch
				if !ok {
					defer wg.Done()
					break
				}

				sink <- value
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(sink)
	}()

	return sink
}

func generateWork(work []int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for _, w := range work {
			ch <- w
		}
	}()

	return ch
}

func main() {
	i1 := generateWork([]int{0, 2, 6, 8})
	i2 := generateWork([]int{1, 3, 5, 7})

	out := fanIn(i1, i2)

	for value := range out {
		fmt.Println("Value:", value)
	}
}
