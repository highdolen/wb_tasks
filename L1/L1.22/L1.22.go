package main

import (
	"fmt"
	"math/big"
)

func main() {
	var a, b string
	fmt.Print("Введите число a: ")
	fmt.Scan(&a)
	fmt.Print("Введите число b: ")

	aBig := new(big.Int)
	bBig := new(big.Int)

	aBig.SetString(a, 10)
	bBig.SetString(b, 10)

	sum := new(big.Int).Add(aBig, bBig)
	sub := new(big.Int).Sub(aBig, bBig)
	mul := new(big.Int).Mul(aBig, bBig)
	div := new(big.Int)
	if bBig.Cmp(big.NewInt(0)) != 0 {
		div.Div(aBig, bBig)
	} else {
		fmt.Println("Ошибка деления на ноль")
	}

	fmt.Println(sum, sub, mul, div)
}
