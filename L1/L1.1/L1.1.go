package main

import "fmt"

type Human struct {
	age     int
	height  int
	weight  int
	country string
	name    string
}

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

type Action struct {
	Human
	action string
}

func main() {
	personActivity := Action{
		Human: Human{
			age:     17,
			height:  187,
			weight:  77,
			country: "Russia",
			name:    "Bob",
		},
		action: "is playing footbal now",
	}

	fmt.Printf("%s %s. He is %v years old. His weight: %v, his height: %v. He is from %s", personActivity.name, personActivity.action, personActivity.age, personActivity.weight, personActivity.height, personActivity.country)

}
