package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 2)
	for t := range ticker.C { //* ticker.C - односторонний канал только на чтение, в котором передается время (раз в 2 секунды он передает время)
		fmt.Println(t)
	}
}
