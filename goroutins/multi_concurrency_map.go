package main

func main() {
	m := make(map[int]int)

	go func() {
		for i := 0; i < 10000; i++ {
			m[0] = i
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			if m[i] < 0 { // типа вызов из мапы
				// some action
			}
		}
	}()
}
