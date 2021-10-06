package supervisor

import (
	"bytes"
	"encoding/json"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/priorityqueue"
	"github.com/foxfurry/go_kitchen/internal/service/cook"
	"github.com/spf13/viper"
	"net/http"
	"sync"
	"time"
)

type ISupervisor interface {
	AddOrder(order dto.Order)
	Initialize()
}

type KitchenSupervisor struct {
	queueMutex sync.Mutex
	queue priorityqueue.IQueue

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

func (s *KitchenSupervisor) AddOrder(order dto.Order) {
	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()
	s.queue.Push(order)
}


func (s *KitchenSupervisor) Initialize(){
	logger.LogMessage("Starting kitchen supervisor")
	s.queue = priorityqueue.NewOrderQueue()
	go s.WatchOrderQueue()
}

func (s *KitchenSupervisor) WatchOrderQueue() {
	for {
		s.queueMutex.Lock()
		for idx := 0; idx < s.queue.Len(); idx++{
			if s.findCook(s.queue.Check(idx).MaxOrderRank) != nil {	// We have an free cook for current priority
				go s.prepareOrder(s.queue.Remove(idx))
			}
		}
		s.queueMutex.Unlock()
		time.Sleep(time.Second * 1) // Every second we iterate through
	}
}

func (s *KitchenSupervisor) getOrderQueue() priorityqueue.IQueue {
	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()
	return s.queue
}

func (s *KitchenSupervisor) prepareOrder(order dto.Order) {
	logger.LogMessageF("Preparing order %v", order.OrderID)

	itemChan := make(chan int, 1)

	for idx, val := range order.Items {
		go s.prepareItem(s.foods[val], idx, itemChan)
	}

	for len(order.Items) > 0 {
		<-itemChan

		order.Items = order.Items[:len(order.Items)-1] // If we received a signal - pop 1 item
	}
}

func (s *KitchenSupervisor) prepareItem(item entity.Food, itemIdx int, itemChan chan<- int){
	var currCook *cook.Cook

	for currCook == nil {
		currCook = s.findCook(item.Complexity)
		time.Sleep(time.Second)		// Every second try to find a cook
	}

	cookID := currCook.GetID()
	logger.LogCookF(cookID, "Preparing item %d", item.ID)

	currCook.Prepare(item, itemIdx, itemChan)

	logger.LogCookF(cookID, "Item %d is ready", item.ID)
}

func (s *KitchenSupervisor) findCook(rankReq int) *cook.Cook{
	s.cooksMutex.Lock()
	defer s.cooksMutex.Unlock()
	for idx, _ := range s.cooks {
		if s.cooks[idx].GetState() == cook.Free && s.cooks[idx].Rank >= rankReq {
			return &s.cooks[idx]
		}
	}
	return nil
}





