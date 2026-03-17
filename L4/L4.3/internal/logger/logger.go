package logger

import "log"

// Logger — асинхронный логгер на базе buffered channel
type Logger struct {
	ch chan string
}

// New - создает логгер, запускает воркер и возвращает готовый экземпляр
func New(buffer int) *Logger {
	l := &Logger{
		ch: make(chan string, buffer),
	}

	go l.worker()

	return l
}

// worker — фоновая горутина, которая читает сообщения из канала и выводит их в stdout
func (l *Logger) worker() {
	for msg := range l.ch {
		log.Println(msg)
	}
}

// Log - пытается неблокирующе записать сообщение в канал логгера
func (l *Logger) Log(msg string) {
	select {
	case l.ch <- msg:
	default:
	}
}
