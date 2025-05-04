// fanin + worker pool
package main

import (
	"fmt"
	"sync"
)

func main() {
	// Создаем входные каналы
	inputChs := make([]chan int, 10)
	for i := range inputChs {
		inputChs[i] = make(chan int)
		go generateNumbers(inputChs[i], i) // Генерация данных
	}

	outputCh := merge(inputChs)
	// Выводим результаты
	for num := range outputCh {
		fmt.Println("Received:", num)
	}
}

// Генерация чисел для примера
func generateNumbers(ch chan int, workerID int) {
	for i := 0; i < 3; i++ {
		ch <- workerID*10 + i
	}
	close(ch)
}

func merge(chs []chan int) chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	poolCh := make(chan chan int)

	go func() {
		for _, ch := range chs {
			poolCh <- ch
		}
		close(poolCh)
	}()

	wg.Add(5)
	for range 5 {
		go func() {
			defer wg.Done()
			for ch := range poolCh {
				for value := range ch {
					out <- value
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
