package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
   1. Сделать структуры Base и Child.
   2. Структура Base должна содержать строковое поле name.
   3. Структура Child должна содержать строковое поле lastName.
   4. Сделать функцию Say у структуры Base, которая распечатывает на экране: Hello, %name!
   5. Пронаследовать Child от Base.
   6. Инициализировать экземпляр b1 Base.
       присвоить name значение Parent
   7. Инициализировать экземпляр c1 Сhild.
       присвоить name значение Child
       присвоить lastName значение Inherited
   8. Вызвать у обоих экземпляров метод Say.
   9. Переопределить метод Say для стурктуры Child, чтобы он выводил на экран: Hello, %lastName %name!
   10. Сделать массив, содержащий b1 и c1.
   11. Вызвать Say у всех элементов массива из шага 10.
   12. Сделать метод NewObject для создания экземпляров Base и Child в зависимости от входного параметра.
   13. Написать юнит тесты для метода NewObject.
   14. Сделать генератор объетов Base и Child такой, чтобы:
       объекты Base создавались в фоновом потоке с задержкой 1 секунда;
       объекты Child создавались в фоновом потоке с задержкой 2 секунды;
       общее время генерации объектов не превышало 11 секунд;
   15. Сделать асинхронный обработчик сгенерированных объектов такой, чтобы:
       метод Say вызывался в порядке генерации объектов;
       не приводил к утечкам памяти;
*/

type Sayer interface {
	Say()
}

type Base struct {
	name string
}

type Child struct {
	Base
	lastName string
}

func NewObject(params []string) (Sayer, error) {
	if len(params) == 1 {
		return Base{name: params[0]}, nil
	} else if len(params) == 2 {
		return Child{Base{params[0]}, params[1]}, nil
	} else {
		return nil, errors.New("some error")
	}

}

func (c Child) Say() {
	fmt.Printf("Hello %s %s\n", c.lastName, c.name)
}

func (b Base) Say() {
	fmt.Printf("Hello %s\n", b.name)
}

func main() {
	wg := &sync.WaitGroup{}
	out := ObjGenerator(wg)
	wg.Add(1)
	go SayObj(wg, out)

	wg.Wait()
}

func ObjGenerator(wg *sync.WaitGroup) <-chan Sayer {
	timeout, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	out := make(chan Sayer)
	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		firstTicker := time.NewTicker(time.Second)
		secondTicker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-timeout.Done():
				close(out)
				return
			case <-firstTicker.C:
				value, err := NewObject([]string{"name"})
				if err != nil {
					return
				}
				out <- value
			case <-secondTicker.C:
				value, err := NewObject([]string{"name", "lastname"})
				if err != nil {
					return
				}
				out <- value
			}
		}
	}()
	return out
}

func SayObj(wg *sync.WaitGroup, in <-chan Sayer) {
	defer wg.Done()
	for value := range in {
		value.Say()
	}
}
