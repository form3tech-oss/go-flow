package api

import (
	"github.com/form3tech/go-logger/log"
	"os"
)

func Start(ch <-chan os.Signal, startedSignal chan bool, stopped chan struct{}) {
	if err := StartServer(ch, startedSignal, stopped); err != nil {
		log.Errorf("Server failed to start: %s", err)
	}
}