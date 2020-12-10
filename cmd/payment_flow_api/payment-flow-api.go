package main

import (
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	api.Configure()

	api.Start(stop, make(chan bool), make(chan struct{}))
}
