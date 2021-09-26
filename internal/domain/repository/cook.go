package repository

import (
	"encoding/json"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
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
		json.Unmarshal(byteValue, &cooksHolder)

		cooks = cooksHolder.Data
	})

	return cooks
}