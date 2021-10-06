package main

import (
	"github.com/foxfurry/go_kitchen/application"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/config"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/profiler"
	"os"
	"os/signal"
	"syscall"
)

func init(){
	config.LoadConfig()
}

func main(){
	profiler.StartProfiler()

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	app := application.CreateApp()
	go app.Start()

	<-sigChan
}
