package main

//создаем структуру(тип данных) c полем msg
type customError struct {
	msg string
}

//реализуем метод Error(), теперь customError удовлетворяет интерфейсу error
func (e *customError) Error() string {
	return e.msg
}

//функция возвращает указатель на customError. data = nil, tab != nil(т.к. определен конкретный тип. под капотом структуры itab, _type = *customError)
func test() *customError {
	// ... do something
	return nil
}

func main() {
	//создаем переменную типа error. Под капотом iface с tab = nil, data = nil
	var err error
	//присваиваем результат функции test. data = nil, tab != nil
	err = test()
	//заходим в данное условие, т.к. tab != nil
	if err != nil {
		//печатаем error
		println("error")
		//выходим из main-горутины
		return
	}
	println("ok")
}
