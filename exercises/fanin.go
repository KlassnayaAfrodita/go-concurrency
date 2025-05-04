//дан канал каналов, нужно слить все в один канал, нужен контекст,
//при отмене которого не будет пушей

package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	mainCh := make(chan chan int)
	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		ch := generate(wg, i)
		mainCh <- ch
	}
	close(mainCh)
	out := mergeChannels(wg, mainCh)

	for value := range out {
		fmt.Println(value)
	}
}

func mergeChannels(ctx context.Context, in <-chan <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ch, ok := <-in:
				if !ok {
					return
				}
				go func() {
					for {
						select {
						case <-ctx.Done():
							return
						case num, ok1 := <-ch:
							if !ok1 {
								return
							}
							select {
							case <-ctx.Done():
								return
							case out <- num:
								return
							}
							out <- num
						}
					}
				}()

			}
		}
	}()

	return out
}

func generate(wg *sync.WaitGroup, num int) <-chan int {
	out := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			out <- num
		}
		close(out)
	}()

	return out
}
