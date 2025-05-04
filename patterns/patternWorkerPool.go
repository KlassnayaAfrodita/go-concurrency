//? why use - потому что это паттерн ограниченной асинхрончины. то есть, когда нужно выполнить работу асинзронно,
//? но ограничены ресурсы

//* when use - когда нужен параллизм/асинхрончина, но ресурсы ограничены

//! solution - использованием фиксированного количества worker для выполнения нескольких задач в очереди.
//! В экосистеме Go мы используем горутины для создания worker и реализации очереди с помощью channel
//! Определенное количество workerов извлечет задачу из очереди и завершит ее, а когда задача будет выполнена,
//! worker продолжит извлекать новую, пока очередь не опустеет.

//* У вас нет неограниченных ресурсов на вашей машине, минимальный размер объекта goroutine составляет 2 КБ,
//* когда вы создаете слишком много goroutine , ваша машина быстро исчерпает память,
//* и CPU будет продолжать обрабатывать задачу, пока не достигнет предела.
//* Используя ограниченный пул рабочих процессов и сохраняя задачу в очереди,
//* мы можем уменьшить нагрузку на CPU и память, поскольку задача будет ждать в очереди,
//* пока рабочий процесс не вытащит задачу.

package main

import (
	"log"
	"sync"
	"time"
)

// type T interface{} // принимаем любой тип

type WorkerPool interface {
	Run(wg *sync.WaitGroup)
	AddTask(task func())
}

type workerPool struct {
	maxWorker         int
	queuedTaskChannel chan func()
}

func NewWorkerPool(maxWorker int) *workerPool {
	return &workerPool{
		maxWorker:         maxWorker,
		queuedTaskChannel: make(chan func()),
	}
}

func (wp *workerPool) Run(wg *sync.WaitGroup) {
	for i := 0; i < wp.maxWorker; i++ {
		go func(workerID int) {
			defer wg.Done()
			for task := range wp.queuedTaskChannel {
				task()
			}
		}(i)
	}
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTaskChannel <- task
}

func main() {
	var wg sync.WaitGroup

	// go func() {
	// 	for {
	// 		log.Printf("[main] Total current goroutine: %d", runtime.NumGoroutine())
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// Start Worker Pool.
	totalWorker := 3
	wp := NewWorkerPool(totalWorker) //* создаем воркер пуд

	wg.Add(totalWorker)
	wp.Run(&wg) //* запускаем воркеров, ждем, когда в канал чтото придет

	type result struct {
		id    int
		value int
	}

	totalTask := 5
	resultC := make(chan result, totalTask) //* создаем канал результатов. буферезацией на сколько задач

	for i := 0; i < totalTask; i++ { //* запускаем таски
		id := i
		wp.AddTask(func() { //? передаем таску (функцию) в очередь пула
			log.Printf("[main] Starting task %d", id)
			time.Sleep(5 * time.Second)
			resultC <- result{id, id * 2}
		})
	}

	for i := 0; i < totalTask; i++ { //* считываем резы.
		res := <-resultC
		log.Printf("[main] Task %d has been finished with result %d:", res.id, res.value)
	}

	go func() {
		wg.Wait()
	}()
}
