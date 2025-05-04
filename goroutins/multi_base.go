package main

import (
	"fmt"
	"time"
)

type printer struct {
}

func (printer) PrintHello() {
	fmt.Println("hello from struct method")
}

func PrintHello() {
	fmt.Println("hello from named func")
}

func main() {
	go func() { // вызов анонимной функции
		fmt.Println("hello from anonymous func")
		time.Sleep(time.Second * 2)
		fmt.Println("hello from anonymous func. Part 2") // не выведется тк ждем больше, чем в main
	}()

	go PrintHello() // вызов известной функции

	var p printer
	go p.PrintHello() // вызов из класса (метода)

	time.Sleep(time.Second)
}
