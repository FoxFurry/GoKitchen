package controller

import (
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/http/httperr"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/foxfurry/go_kitchen/internal/service/supervisor"
	"github.com/gin-gonic/gin"
	"net/http"
)

const CurrentCaller = "Kitchen Controller"

type IController interface {
	menu(c *gin.Context)
	order(c *gin.Context)
	RegisterKitchenRoutes(c *gin.Engine)
	Initialize()
}

type KitchenController struct {
	super supervisor.ISupervisor
}

func NewKitchenController() IController {
	return &KitchenController{
		super: supervisor.NewKitchenSupervisor(),
	}
}

func (ctrl *KitchenController) Initialize(){
	ctrl.super.Initialize()
}

func (ctrl *KitchenController) menu(c *gin.Context){
	var response dto.Menu

	response.Items = repository.GetFoods()
	response.ItemsCount = len(response.Items)
	logger.LogMessageF("Menu request was fulfilled: %d items available", response.ItemsCount)
	c.JSON(http.StatusOK, response)
}

func (ctrl *KitchenController) order(c *gin.Context){
	var currentOrder dto.Order

	if err := c.ShouldBindJSON(&currentOrder); err != nil {
		httperr.HandleErr(CurrentCaller, err, c)
		return
	}


	logger.LogMessageF("Got a new order: %v", currentOrder.Items)
	ctrl.super.AddOrder(currentOrder)

	return
}

/*
logger.LogMessageF("Order %v completed", currentOrder.OrderID)

	resp := dto.Distribution{}
	resp.TableID = currentOrder.TableID

	jsonBody, err := json.Marshal(resp)
	if err != nil {
		log.Panic(err)
	}
	contentType := "application/json"

	http.Post(viper.GetString("dining_host") + "/distribution", contentType, bytes.NewReader(jsonBody))

 */
