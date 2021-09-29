package main

import (
	"github.com/foxfurry/go_kitchen/application"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/config"
	"os"
	"os/signal"
	"syscall"
)

func init(){
	config.LoadConfig()
}

func main(){
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	app := application.CreateApp()
	go app.Start()

	<-sigChan
}
