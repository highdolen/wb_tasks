package main

import "fmt"

//создание стуктуры
type Human struct {
	age     int
	height  int
	weight  int
	country string
	name    string
}

//создание методов структуры Human
func (h *Human) GetAge() int {
	return h.age
}

func (h *Human) GetHeight() int {
	return h.height
}

func (h *Human) GetWeight() int {
	return h.weight
}

func (h *Human) GetCountry() string {
	return h.country
}

func (h *Human) GetName() string {
	return h.name
}

//создание структуры Action с использованием композиции embedded struct(делаем так, чтобы структура Action имела все все методы Human)
type Action struct {
	Human
	action string
}

func main() {
	personActivity := Action{ //заполняем поля структуры Action
		Human: Human{ //заполняем поля структуры Human
			age:     17,
			height:  187,
			weight:  77,
			country: "Russia",
			name:    "Bob",
		},
		action: "is playing footbal now",
	}
	//проверям, имеет ли структура Action все поля структуры Human
	fmt.Printf("%s %s. He is %v years old. His weight: %v, his height: %v. He is from %s", personActivity.name, personActivity.action, personActivity.age, personActivity.weight, personActivity.height, personActivity.country)
	//на выходе получаем нужный результат, который подтверждает, что структура Action имеет доступ к методам Human
}
