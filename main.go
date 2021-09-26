package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	app := application.CreateApp()
	go app.Start()

	<-sigChan

	app.Shutdown(ctx)
	cancel()
}
