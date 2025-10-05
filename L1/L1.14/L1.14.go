package main

import (
	"fmt"
	"reflect"
)

func detect(value any) string {

	//если переменная равна nil, печатаем nil
	if value == nil {
		return "nil"
	}

	//проверям тип переменных с помощью конструкции switch-case, также используется оператор типа
	switch value.(type) {
	//если тип переменной int, то возвращаем int
	case int:
		return "int"
	//если тип переменной string, то возвращаем string
	case string:
		return "string"
	//если тип переменной bool, то возвращаем bool
	case bool:
		return "bool"
	//так как осталось проверить только chan, можем написать default
	default:
		//если значение является каким либо типом chan(chan int, chan bool и т.д.), тогда сообщаем о том
		//что это переменная chan
		if reflect.ValueOf(value).Kind() == reflect.Chan {
			return "chan"
		} else {
			//для всех остальных типов, говорим о том, что это неизвестный тип
			return "no name type"
		}
	}

}

func main() {
	//печатаем информацию о том, каким типом является переменная
	fmt.Println(detect(52))
	fmt.Println(detect("word"))
	fmt.Println(detect(false))
	fmt.Println(detect(nil))
	fmt.Println(detect(make(chan struct{})))
}
