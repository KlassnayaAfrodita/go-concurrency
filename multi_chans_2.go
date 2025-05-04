package main

import "fmt"

func main() {
	first := make(chan int)
	prev := first
	const goroutinesCount = 1e2

	for i := 0; i < goroutinesCount; i++ {
		next := make(chan int)

		go func(prev chan int) {
			number := <-prev //! prev пустой и планировщик пропускает этот шаг, идет дальше, пока не положим чтонибудь, поэтому мы ждем, когда в него чтото положат
			next <- number
		}(prev)

		prev = next
	}
	first <- 42 //! здесь мы кладем в prev
	fmt.Println(<-prev)
}

//todo  те мы запускаем 100 горутин и они ждут, пока в канал не придет какоенибудь значение
//todo  затем приходит значение и оно передается через каналы к горутинам
