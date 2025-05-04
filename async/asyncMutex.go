package main

import (
	"fmt"
	"sync"
	"time"
)

// * критической секции не должно быть много. добавляем только туда, где есть конкурентный (эксклюзивный) доступ
type Counter struct {
	mu sync.Mutex
	c  map[string]int
}

func (c *Counter) Inc(key string) {
	//* когда одна горутина залочена, в метод заходит другая, видит, что горутина залочена и ждет
	//* таким образом горутины выстраиваются в очередь

	c.mu.Lock() //? исправили с помощью мьютекса
	c.c[key]++  //* критическая секция
	//! как мы помним, мапа - непотокобезопасный тип. те много горутин обращаются и происходят пересечения
	c.mu.Unlock()
}

func (c *Counter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c[key] //? критическая секция
}

type RWCounter struct {
	mu sync.RWMutex
	c  map[string]int
}

func (c *RWCounter) CountMe() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c
}

func (c *RWCounter) CountMeAgain() map[string]int {
	c.mu.RLock() //! блокируем только на чтение
	defer c.mu.RUnlock()
	return c.c
}

func main() {
	key := "test"

	c := Counter{c: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.Inc(key)
	}

	time.Sleep(3 * time.Second)
	fmt.Println(c.Value(key))
}
