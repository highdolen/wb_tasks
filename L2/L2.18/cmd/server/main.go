package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Context(context.Background())
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Errorf("Ошибка в загрузке конфига:%v", err)
	}
}
