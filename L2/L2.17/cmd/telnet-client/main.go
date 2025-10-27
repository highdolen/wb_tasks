package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"secondBlock/internal/client"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout timeout] host port\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Defaults:\n")
		fmt.Fprintf(os.Stderr, "  --timeout 10s\n")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	config := client.Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}

	telnetClient := client.New(config)
	if err := telnetClient.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
