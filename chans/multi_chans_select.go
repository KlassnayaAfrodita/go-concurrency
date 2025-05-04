package main

import "fmt"

func main() {
	c1 := make(chan int, 5)
	c2 := make(chan int, 5)

	for i := 1; i <= 3; i++ {
		c1 <- i
		c2 <- -i
	}
	// close(c1)
	// close(c2) //! когда мы читаем из закрытого канала, нам дается дефолтное значение для типа данных (для int - 0)

	for i := 1; i < 10; i++ {
		select { // одновременное ожидание чтения/записи
		//* если оба case удовлетворены, то поведение селекта непредсказуемо
		case x := <-c1:
			fmt.Println("c1:", x)
		case x := <-c2: // ожидаем чтение с с2
			fmt.Println("c2:", x)
		default:
			fmt.Println("no data!")
		}
	}
}
