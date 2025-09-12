package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Получаем текущее время через NTP
	currentTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		// Ошибка выводится в STDERR
		fmt.Fprintln(os.Stderr, "Ошибка получения времени:", err)
		os.Exit(1)
	}

	// Выводим время в читаемом формате
	fmt.Println(currentTime.Format(time.RFC1123))
}