package main

//импортируем нужные библиотеки
import (
	"fmt"
	"os"
)

func Foo() error {
	// под капотом iface(т.к. error это интерфейс с методом). data = nil, tab != nil(определен конкретный тип *os.PathError)
	var err *os.PathError = nil
	return err
}

func main() {
	//инициализируем переменную err(возвращаем тип error(под капотом iface))
	err := Foo()
	fmt.Println(err)        //nil(выводим именно значение(data))
	fmt.Println(err == nil) //false(value = nil, tab != nil, т.к. известен тип). Для интерфейса вывело бы true при условии data = nil, tab = nil
}
