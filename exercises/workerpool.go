// Реализуйте обработку 1000 URL-запросов с использованием пула из 5 воркеров.
// Каждый воркер должен отправлять HTTP GET-запрос и возвращать статус ответа.
// Ограничьте количество одновременных запросов через пул горутин.
package main

import (
	"fmt"
	"sync"
)

func main() {
	urls := make([]string, 10)
	urlsCh := generateURL(urls)
	outputCh := make(chan int)
	wg := sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			//defer close(outputCh)
			defer wg.Done()
			for url := range urlsCh {
				value := sendRequest(url)
				outputCh <- value
			}

			return
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

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
