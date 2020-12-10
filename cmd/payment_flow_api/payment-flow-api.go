package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	api.Configure()

	api.Start(stop, make(chan bool), make(chan struct{}))
}
