package controller

import "github.com/gin-gonic/gin"

func (ctrl *KitchenController) RegisterKitchenRoutes(c *gin.Engine){
	c.GET("/menu", ctrl.menu)
	c.POST("/order", ctrl.order)
}

