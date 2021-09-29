package application

import (
	"github.com/foxfurry/go_kitchen/internal/http/controller"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type IApp interface {
	Start()
}

type kitchenApp struct {
	router *gin.Engine
}

func CreateApp() IApp {
	app := kitchenApp{
		router: gin.Default(),
	}

	ctrl := controller.NewKitchenController()
	ctrl.RegisterKitchenRoutes(app.router)

	return &app
}

func (d *kitchenApp) initialize(){

}

func (d *kitchenApp) Start() {
	d.initialize()
	logger.LogMessage("Starting kitchen server")

	if err := d.router.Run(viper.GetString("kitchen_host")); err != http.ErrServerClosed {
		logger.LogPanicf("Unexpected error while running server: %v", err)
	}
}