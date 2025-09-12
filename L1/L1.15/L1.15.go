package main

var justString string

func someFunc() {
	v := createHugeString(1 << 10)
	justString = v[:100]                 //неправильно
	justString = string([]rune(v[:100])) //правильно, создаем новый слайс и работаем с новыми данными, также можно вместо []rune написать []byte
}

func main() {
	someFunc()
}
