package repository

import (
	"encoding/json"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"io/ioutil"
	"os"
	"sync"
)

var (
	foods []entity.Food
	onceFood sync.Once
)

func GetFoods() []entity.Food{
	onceFood.Do(func(){
		foodsHolder := struct {
			Data []entity.Food `json:"foods"`
		}{}

		jsonFile, _ := os.Open("./config/foods.json")
		byteValue, _ := ioutil.ReadAll(jsonFile)

		if err := json.Unmarshal(byteValue, &foodsHolder); err != nil {
			logger.LogPanicf("Could not unmarshal foods config: %v", err)
		}

		foods = foodsHolder.Data
	})

	return foods
}
