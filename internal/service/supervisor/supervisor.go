package supervisor

import (
	"bytes"
	"encoding/json"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/foxfurry/go_kitchen/internal/service/cook"
	"github.com/spf13/viper"
	"net/http"
	"sync"
	"time"
)

type ISupervisor interface {
	PrepareOrder(order dto.Order)
}

type KitchenSupervisor struct {
	cooksMutex sync.Mutex
	cooks []cook.Cook

	foodsMutex sync.Mutex
	foods []entity.Food
}

func NewKitchenSupervisor() ISupervisor {
	return &KitchenSupervisor{
		cooks: cook.EntityToService(repository.GetCooks()),
		foods: repository.GetFoods(),
	}
}

func (s *KitchenSupervisor) PrepareOrder(order dto.Order) {
	s.prep(order)
	logger.LogSuperF("Order %v completed", order.OrderID)

	resp := dto.Distribution{}
	resp.TableID = order.TableID

	jsonBody, err := json.Marshal(resp)
	if err != nil {
		logger.LogPanic(err.Error())
	}
	contentType := "application/json"

	http.Post(viper.GetString("dining_host") + "/distribution", contentType, bytes.NewReader(jsonBody))
}

func (s *KitchenSupervisor) prep(order dto.Order) {
	logger.LogSuperF("Got order %v", order.OrderID)

	itemChan := make(chan int, 1)

	for idx, val := range order.Items {
		go s.PrepareItem(s.foods[val], idx, itemChan)
	}

	for len(order.Items) > 0 {
		<-itemChan

		order.Items = order.Items[:len(order.Items)-1]
	}
}

func (s *KitchenSupervisor) FindCook() *cook.Cook{
	s.cooksMutex.Lock()
	defer s.cooksMutex.Unlock()
	for idx, _ := range s.cooks {
		if s.cooks[idx].GetState() == cook.Free {
			s.cooks[idx].SetState(cook.Busy)
			return &s.cooks[idx]
		}
	}
	return nil
}

func (s *KitchenSupervisor) PrepareItem(item entity.Food, itemIdx int, itemChan chan<- int){
	var currCook *cook.Cook

	for currCook == nil {
		currCook = s.FindCook()
		time.Sleep(time.Second)
	}

	cookID := currCook.GetID()
	logger.LogCookF(cookID, "Preparing item %d", item.ID)

	currCook.Prepare(item, itemIdx, itemChan)

	logger.LogCookF(cookID, "Item %d is ready", item.ID)
}




