//Реализуйте обработку 1000 URL-запросов с использованием семафор из 5 воркеров.
//Каждый воркер должен отправлять HTTP GET-запрос и возвращать статус ответа.
//Ограничьте количество одновременных запросов через пул горутин.

package main

import (
	"fmt"
	"sync"
)

func main() {
	urls := make([]string, 10)
	urlsCh := generateURL(urls)
	wg := sync.WaitGroup{}
	sema := make(chan struct{}, 5)
	outputCh := make(chan int, len(urls)) // можно буферизовать, чтобы избежать блокировок

	for url := range urlsCh {
		wg.Add(1)
		sema <- struct{}{}
		go func(url string) {
			defer wg.Done()
			defer func() { <-sema }()

			value := sendRequest(url)
			outputCh <- value
		}(url)
	}

	// Закрываем outputCh после завершения всех запросов.
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	// Читаем и выводим результаты в main – это гарантирует, что программа не завершится раньше.
	for value := range outputCh {
		fmt.Println(value)
	}
}

func sendRequest(url string) int {
	_ = url
	return 200
}

func generateURL(urls []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, url := range urls {
			out <- url
		}
		close(out)
	}()
	return out
}
