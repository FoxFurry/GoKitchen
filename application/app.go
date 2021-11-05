package application

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/http/controller"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"time"
)

type IApp interface {
	Start()
	Shutdown(ctx context.Context)
}

type kitchenApp struct {
	server *http.Server
}

func Create(ctx context.Context) IApp {
	appHandler := gin.New()

	ctrl := controller.NewKitchenController()
	ctrl.RegisterKitchenRoutes(appHandler)

	app := kitchenApp{
		server: &http.Server{
			Addr:    viper.GetString("kitchen_host"),
			Handler: appHandler,
		},
	}

	if viper.GetInt("version") == 2 {
		app.registerKitchen()
	}

	ctrl.Initialize(ctx)

	return &app
}

func (d *kitchenApp) Start() {
	logger.LogMessage("Starting kitchen server")

	if err := d.server.ListenAndServe(); err != http.ErrServerClosed {
		logger.LogPanicF("Unexpected error while running server: %v", err)
	}
}

func (d *kitchenApp) Shutdown(ctx context.Context) {
	if err := d.server.Shutdown(ctx); err != nil {
		logger.LogPanicF("Unexpected error while closing server: %v", err)
	}
	logger.LogMessage("Server terminated successfully")
}

func (d *kitchenApp) registerKitchen() {
	aggregatorHost := viper.GetString("aggregator_host")
	deliveryHost := viper.GetString("delivery_host")

	logger.LogMessageF("Trying to reach aggregator server on: %v", aggregatorHost)
	waitConnection(aggregatorHost)
	logger.LogMessageF("Reached aggregator server!")

	logger.LogMessageF("Trying to reach delivery server on: %v", deliveryHost)
	waitConnection(deliveryHost)
	logger.LogMessageF("Reached delivery server!")

	items := repository.GetFoods()
	registerData := dto.RestaurantRegister{
		RestaurantID: 1,
		Name:         viper.GetString("restaurant_name"),
		Address:      "http://localhost" + viper.GetString("kitchen_host"),
		MenuItems:    len(items),
		Menu:         items,
		Rating:       10,
	}

	contentType := "application/json"

	for {
		jsonBody, err := json.Marshal(registerData)
		if err != nil {
			log.Panic(err)
		}

		resp, err := http.Post("http://" + viper.GetString("aggregator_host")+"/register", contentType, bytes.NewReader(jsonBody))
		if err != nil {
			logger.LogWarningF("Could not register restaurant: %s. Restaurant is working autonomously")
			break
		}

		if resp.StatusCode == 200 {
			logger.LogMessageF("Successfully registered restaurant %s", registerData.Name)
			break
		}

		logger.LogWarningF("Could not register restaurant: %s. Trying with different ID: %d", resp.Status, registerData.RestaurantID+1)
		registerData.RestaurantID++

		time.Sleep(1)
	}
}

func waitConnection(host string) {
	tryCount := 1
	dialTick := time.Tick(time.Second * time.Duration(viper.GetInt("dial_timeout")))

	for {
		select {
		case <-dialTick:
			conn, err := net.Dial("tcp", host)

			if err == nil {
				conn.Close()
				return
			}

			logger.LogWarningF("Could not reach the server: %v. Retrying %d", err, tryCount)
			tryCount++
		}
	}
}
