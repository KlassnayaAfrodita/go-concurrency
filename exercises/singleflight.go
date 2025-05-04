//задача: есть платная ручка requestThatCostsMoney, надо сделать так, чтобы при
//параллельном запросе в нее, шел только один запрос

// рассуждения: нужно сделать так, чтобы проходила только одна горутинка, остальные стопились,
// аналог - при отмене родительского контекста, отменяется и дочерний
// надо отправлять всем горутинам сигналам
package main

import "sync"

func main() {
	mo := MoneyOptimization{}
	userID := "vasya"

	go func() {
		mo.RequestThatCostsMoney(userID)
	}()

	go func() {
		mo.RequestThatCostsMoney(userID)
	}()

	go func() {
		mo.RequestThatCostsMoney(userID)
	}()

	go func() {
		mo.RequestThatCostsMoney(userID)
	}()

	// тут еще может быть множество параллельно запускаемых запросов
}

func requestThatCostsMoney(userID string) (response string) {
	_ = userID
	// Функция, которая стоит денег
	return "some"
}

type Some struct {
	response string
	ch       chan struct{}
}

type MoneyOptimization struct {
	// мы создаем мапу, чтобы както определять, когда у нас одинаковые userID
	cache map[string]*Some // здесь ключ userID
	mu    sync.Mutex
}

// сейчас разбираем горутину, которая дошла первая
// она лочит мьютекс, не находит в кэше ничего
func (o *MoneyOptimization) RequestThatCostsMoney(userID string) (response string) {
	// тут пиши обертку над функцией requestThatCostsMoney
	o.mu.Lock()                          // лочим для похода в мапу
	if some, ok := o.cache[userID]; ok { // смотрим, есть ли уже ответ в мапе
		o.mu.Unlock()

		<-some.ch // сюда попали все горутины, кроме первой, они ждут, пока первая сделает работу
		return some.response
	}

	o.cache[userID] = &Some{ch: make(chan struct{})}
	o.mu.Unlock()

	o.mu.Lock()
	resp := requestThatCostsMoney(userID)
	o.mu.Unlock()

	o.mu.Lock()
	o.cache[userID].response = resp
	o.mu.Unlock()

	close(o.cache[userID].ch)

	o.mu.Lock()
	delete(o.cache, userID)
	o.mu.Unlock()

	return resp
}
