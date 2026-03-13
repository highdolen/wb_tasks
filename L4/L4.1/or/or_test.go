package or

import (
	"testing"
	"time"
)

func Sig(after time.Duration) <-chan interface{} {
	//создаем канал для сигнала о завершении работы
	c := make(chan interface{})
	//горутина, которая закрывает канал по истечении времени
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	//возвращаем канал
	return c
}

// TestOrMultiplyChannel - тест нескольких каналовы
func TestOrMultiplyChannels(t *testing.T) {

	start := time.Now()

	<-Or(
		Sig(time.Second),
		Sig(2*time.Second),
		Sig(2*time.Minute),
		Sig(time.Second),
	)

	duration := time.Since(start)

	t.Log(duration)

	if duration > 1300*time.Millisecond {
		t.Fatal("or did not return after fastest channel")
	}
}

// TestOrSingleChannel - тест одного канала
func TestOrSingleChannel(t *testing.T) {

	start := time.Now()

	<-Or(
		Sig(time.Second),
	)

	duration := time.Since(start)

	t.Log(duration)

	if duration > 1300*time.Millisecond {
		t.Fatal("or did not return after fastest channel")
	}
}

// TestOrZeroChannels - тест без передачи каналов
func TestOrZeroChannels(t *testing.T) {

	if Or() != nil {
		t.Fatal("expected nil when no channels are passed")
	}

}
