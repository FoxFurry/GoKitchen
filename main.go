package main

import (
	"context"
	"github.com/foxfurry/go_kitchen/application"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/config"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/gui"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/profiler"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	config.LoadConfig()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	kitchenUI := gui.NewKitchenCUI()

	gui.AppMode = gui.CMDMode
	ok := kitchenUI.Create()
	if ok {
		gui.AppMode = gui.CUIMode
	}

	app := application.Create(ctx)
	go app.Start()
	go profiler.Start(ctx)
	go kitchenUI.Start(ctx, cancel)

	for{
		select {
		case data := <-logger.LogChannel:
			kitchenUI.AddLog(data)
		case <-ctx.Done():
			return
		}
	}
}
