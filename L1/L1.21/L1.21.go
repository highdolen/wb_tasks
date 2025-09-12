package main

import "fmt"

type NewCup interface {
	SystemPoints(msg string)
}

type OldCup struct {
	points int
}

func (o *OldCup) OldSystemPoints(msg string) int {
	if msg == "draw" {
		o.points++
	} else if msg == "win" {
		o.points += 3
	} else if msg == "lose" {
		o.points += 0
	} else {
		fmt.Println("Unknown result")
		return 0
	}
	return o.points
}

type AdapterCup struct {
	Old *OldCup
}

func (a *AdapterCup) SystemPoints(msg string) {
	fmt.Println("Adapter in work")
	fmt.Printf("Points %v", a.Old.OldSystemPoints(msg))
}

func main() {
	var cup NewCup
	cup = &AdapterCup{
		Old: &OldCup{
			points: 5,
		},
	}

	fmt.Println("Введите результат:")
	var result string
	fmt.Scan(&result)
	cup.SystemPoints(result)
}
