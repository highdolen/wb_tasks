package main

//импортируем нужные библиотеки
import (
	"fmt"
	"os"
	"time"

	//импортируем библиотеку ntp с github
	"github.com/beevik/ntp"
)

func main() {
	// получаем текущее время через NTP
	currentTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		// ошибка выводится в STDERR
		fmt.Fprintln(os.Stderr, "Ошибка получения времени:", err)
		os.Exit(1)
	}

	// выводим время в читаемом формате
	fmt.Println(currentTime.Format(time.RFC1123))
}
