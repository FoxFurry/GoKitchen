package supervisor

import (
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/foxfurry/go_kitchen/internal/service/cook"
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
		cooks: cook.CookEntityToService(repository.GetCooks()),
		foods: repository.GetFoods(),
	}
}

func (s *KitchenSupervisor) PrepareOrder(order dto.Order) {
	logger.LogMessageF("Preparing order %v", order.OrderID)

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
	logger.LogMessageF("Found cook")

	currCook.Prepare(item, itemIdx, itemChan)

	logger.LogMessageF("Item %d is ready", item.ID)
}




