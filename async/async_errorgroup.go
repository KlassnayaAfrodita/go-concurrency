package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	wg, wgCtx := errgroup.WithContext()

	for i := 0; i < 3; i++ {
		i := i
		wg.Go(func() error {
			for j := 0; j < 10; j++ {
				if wgCtx.Err() != nil {
					return wgCtx.Err()
				}
				// if i == 1 && j == 3 {
				// 	return fmt.Errorf("some error")
				// }

				fmt.Println("i ", i, "j ", j)
				// runtime.Gosched()
			}
			return nil
		})
	}
	err := wg.Wait()
	fmt.Println(err)
}
