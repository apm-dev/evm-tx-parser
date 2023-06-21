package common

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func WaitForSignal() {
	var stop = make(chan struct{})
	go func() {
		var sig = make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sig)
		<-sig
		log.Println("got interrupt, shutting down...")
		stop <- struct{}{}
	}()
	<-stop
}
