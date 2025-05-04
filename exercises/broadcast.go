package main

import (
	"fmt"
	"sync"
)

type Broadcast struct {
	mu    sync.Mutex
	chans map[chan string]struct{}
}

func NewBroadcast() *Broadcast {
	return &Broadcast{
		chans: make(map[chan string]struct{}),
	}
}

func (b *Broadcast) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan string, 10)
	b.chans[ch] = struct{}{}
	return ch
}

func (b *Broadcast) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.chans, ch)
	close(ch)
}

func (b *Broadcast) Publish(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for ch := range b.chans {
		select {
		case ch <- msg:
		default: // если канал полный - пропускаем
		}
	}
}

// Пример использования
func main() {
	bc := NewBroadcast()

	// Подписчики
	sub1 := bc.Subscribe()
	sub2 := bc.Subscribe()

	go func() {
		for msg := range sub1 {
			fmt.Println("Sub1:", msg)
		}
	}()

	go func() {
		for msg := range sub2 {
			fmt.Println("Sub2:", msg)
		}
	}()

	bc.Publish("Hello!")
	bc.Publish("Broadcast message")

	bc.Unsubscribe(sub1)
	bc.Publish("Last message")
}
