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
	apparatuses    []entity.Apparatus
	onceApparatus sync.Once
)

func GetApparatuses() []entity.Apparatus{
	onceApparatus.Do(func(){
		apparatusesHolder := struct {
			Data []entity.Apparatus `json:"apparatuses"`
		}{}

		jsonFile, _ := os.Open("./config/apparatuses.json")
		byteValue, _ := ioutil.ReadAll(jsonFile)

		if err := json.Unmarshal(byteValue, &apparatusesHolder); err != nil {
			logger.LogPanicF("Could not unmarshal apparatuses config: %v", err)
		}

		apparatuses = apparatusesHolder.Data
	})

	return apparatuses
}
