package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type TelnetClient struct {
	config Config
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewTelnetClient(config Config) *TelnetClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &TelnetClient{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *TelnetClient) connect() error {
	address := net.JoinHostPort(c.config.Host, c.config.Port)

	ctx, cancel := context.WithTimeout(c.ctx, c.config.Timeout)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	c.conn = conn
	log.Printf("Connected to %s\n", address)
	return nil
}

func (c *TelnetClient) Run() error {
	if err := c.connect(); err != nil {
		return err
	}
	defer c.cleanup()

	// Обработка сигналов для graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Каналы для синхронизации горутин
	errCh := make(chan error, 2)

	// Горутина для чтения из сокета и вывода в STDOUT
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		errCh <- c.readFromServer()
	}()

	// Горутина для чтения из STDIN и записи в сокет
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		errCh <- c.writeToServer()
	}()

	// Горутина для обработки сигналов
	go func() {
		select {
		case <-sigCh:
			log.Println("\nReceived interrupt signal, closing connection...")
			c.cancel()
		case <-c.ctx.Done():
		}
	}()

	// Ожидание завершения первой горутины с ошибкой
	err := <-errCh
	c.cancel()

	// Ожидание завершения всех горутин
	c.wg.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (c *TelnetClient) readFromServer() error {
	reader := bufio.NewReader(c.conn)
	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			// Чтение данных из сокета
			data := make([]byte, 1024)
			n, err := reader.Read(data)
			if err != nil {
				if err == io.EOF {
					log.Println("Server closed the connection")
					return nil
				}
				if c.ctx.Err() != nil {
					return nil // Контекст отменен, это нормальное завершение
				}
				return fmt.Errorf("error reading from server: %w", err)
			}

			if n > 0 {
				// Вывод полученных данных в STDOUT
				fmt.Print(string(data[:n]))
			}
		}
	}
}

func (c *TelnetClient) writeToServer() error {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(c.conn)

	for {
		select {
		case <-c.ctx.Done():
			return nil
		default:
			// Чтение из STDIN
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("EOF received, closing connection...")
					return nil
				}
				return fmt.Errorf("error reading from stdin: %w", err)
			}

			// Запись в сокет
			if _, err := writer.WriteString(line); err != nil {
				return fmt.Errorf("error writing to server: %w", err)
			}
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("error flushing to server: %w", err)
			}
		}
	}
}

func (c *TelnetClient) cleanup() {
	if c.conn != nil {
		c.conn.Close()
	}
	c.cancel()
}
