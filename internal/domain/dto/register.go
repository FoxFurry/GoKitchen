package dto

import "github.com/foxfurry/go_kitchen/internal/domain/entity"

type RestaurantRegister struct {
	RestaurantID int `json:"restaurant_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	MenuItems int `json:"menu_items"`
	Menu []entity.Food `json:"menu"`
	Rating float32 `json:"rating"`
}

