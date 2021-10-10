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
	cooks    []entity.Cook
	onceCook sync.Once
)

func GetCooks() []entity.Cook{
	onceCook.Do(func(){
		cooksHolder := struct {
			Data []entity.Cook `json:"cooks"`
		}{}

		jsonFile, _ := os.Open("./config/cooks.json")
		byteValue, _ := ioutil.ReadAll(jsonFile)

		if err := json.Unmarshal(byteValue, &cooksHolder); err != nil {
			logger.LogPanicF("Could not unmarshal cooks config: %v", err)
		}

		cooks = cooksHolder.Data
	})

	return cooks
}