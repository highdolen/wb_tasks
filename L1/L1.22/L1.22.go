package main

//импортируем нужные библиотеки
import (
	"fmt"
	"math/big"
)

func main() {
	//инициализируем две переменные для ввода
	var a, b string
	fmt.Print("Введите число a: ")
	fmt.Scan(&a)
	fmt.Print("Введите число b: ")

	//создаем переменные типа *big.Int для хранения больших чисел
	aBig := new(big.Int)
	bBig := new(big.Int)

	//преобразуем введенные строки чисел в большие числа(СС - 10)
	aBig.SetString(a, 10)
	bBig.SetString(b, 10)

	//суммируем числа
	sum := new(big.Int).Add(aBig, bBig)
	//вычитаем числа
	sub := new(big.Int).Sub(aBig, bBig)
	//умножаем числа
	mul := new(big.Int).Mul(aBig, bBig)
	//инициализируем переменную, где будет храниться результат деления чисел
	div := new(big.Int)
	//если число b(делитель) не равен нулю, тогда выполняем деление
	if bBig.Cmp(big.NewInt(0)) != 0 {
		div.Div(aBig, bBig) //выполняем деление
		//в противном случае печатаем сообщение об ошибке
	} else {
		fmt.Println("Ошибка деления на ноль")
	}

	//печатаем результаты
	fmt.Println(sum, sub, mul, div)
}
