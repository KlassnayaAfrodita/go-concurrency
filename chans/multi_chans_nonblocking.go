package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	c := make(chan struct{})
	go func() {
		fastTicker := time.NewTicker(time.Second * 5 / 2)
		for _ = range fastTicker.C {

			//! неблокирующая запись. если можем записать - пишем, иначе уходим в дефолт
			//! так же делается неблокирующие чтение
			select {
			case c <- struct{}{}: // пишем в с пустую структуру
			default:
				fmt.Println("skip fast job")
			}

		}
	}()

	time.Sleep(time.Second / 2)
	longTicker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-longTicker.C:
			fmt.Println("-----------------long job", time.Since(start))
			time.Sleep(time.Second * 4) // имитируем длинную работу
		case <-c: // читаем с
			fmt.Println("-----------------fast job", time.Since(start))
		}
	}
}

//* сначала мы запускаем горутину и в ней запускаем цикл по тикеру (те раз в 2 секунды он исполняется).
//* дальше делаем неблокирующую запись
//* после запускаем другой тикер. если можем прочитать из него, то делаем длинную работу, если в основном канале есть чтото, то делаем короткую
