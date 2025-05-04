package main

import (
	"fmt"
	"sync"
	"time"
)

// type Task struct {
// 	taskId int
// 	taskDuration int
// 	task func()
// }

type result struct {
	id     int
	result int
}

type TestWorkerPool interface {
	Run(wg *sync.WaitGroup)
	AddTask(task func())
}

type testWorkerPool struct {
	maxWorker         int
	queuedTaskChannel chan func()
}

func (wp *testWorkerPool) testRun(wg *sync.WaitGroup) {
	for i := 0; i < wp.maxWorker; i++ {
		go func(workerID int) {
			defer wg.Done()
			for task := range wp.queuedTaskChannel {
				task()
			}
		}(i)
	}
}

func (wp *testWorkerPool) testAddTask(task func()) {
	wp.queuedTaskChannel <- task
}

func testNewWorkerPool(totalWorkers int) *testWorkerPool {
	return &testWorkerPool{
		maxWorker:         totalWorkers,
		queuedTaskChannel: make(chan func(), totalWorkers),
	}
}

func main() {
	var wg sync.WaitGroup

	totalWorkers := 3
	pool := testNewWorkerPool(totalWorkers)

	wg.Add(3)
	pool.testRun(&wg)

	totalTask := 15

	resultTask := make(chan result, totalTask)

	for i := 0; i < totalTask; i++ {
		pool.testAddTask(func() {
			fmt.Printf("Starting task %d\n", i)
			time.Sleep(1 * time.Second)
			resultTask <- result{
				id:     i,
				result: i * i,
			}
		})
	}

	close(resultTask) //! чтобы не словить дедлок

	for task := range resultTask {
		fmt.Printf("Task %d has been finished with result %d\n", task.id, task.result)
	}

	go func() {
		wg.Wait()
	}()
}
