//? why use - чтобы избежать утечек горутин (когда больше не работает, но запущена)

//* when use - когда надо остановить выполнение горутины

//! solution - done channel

done := make(chan struct{}) // нулевая структура ничего не весит
defer close(done)

s.Run(done)

//! процесс который запустил горутину должен ее же и завершить 

func (s *Service) Run(done <- chan struct{}) <- chan struct{}{
	finished := make(chan struct{})
	// работа
	return finished
}