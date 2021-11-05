package main

import (
	"context"
	"github.com/foxfurry/go_kitchen/application"
	"github.com/foxfurry/go_kitchen/internal/gui"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/config"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/profiler"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	config.LoadConfig()
}

var uiMode = false

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := application.Create(ctx)
	ui := gui.NewKitchenCUI()
	uiMode = ui.Create()

	go app.Start()
	go profiler.Start(ctx)

	if uiMode {
		ui.Start(ctx, cancel)
	}else{
		<-sigs
		cancel()
	}
}

