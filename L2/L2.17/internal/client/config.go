package client

import "time"

type Config struct {
	Host    string
	Port    string
	Timeout time.Duration
}
