package supervisor

import (
	"context"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/domain/repository"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/priorityqueue"
	"github.com/foxfurry/go_kitchen/internal/service/cook"
	"github.com/foxfurry/go_kitchen/internal/service/supervisor/gateway"
	"sync"
	"time"
)

type ISupervisor interface {
	AddOrder(order dto.Order)
	Initialize(ctx context.Context)
}

type KitchenSupervisor struct {
	queueMutex sync.Mutex
	queue      priorityqueue.IQueue

	cooksMutex sync.Mutex
	cooks      []cook.Cook

	foodsMutex sync.Mutex
	foods      []entity.Food
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

func (s *KitchenSupervisor) Initialize(ctx context.Context) {
	logger.LogMessage("Starting kitchen supervisor")
	s.queue = priorityqueue.NewOrderQueue()
	go s.WatchOrderQueue(ctx)
}

func (s *KitchenSupervisor) WatchOrderQueue(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			s.queueMutex.Lock()
			for idx := 0; idx < s.queue.Len(); idx++ {
				if s.findCook(s.queue.Check(idx).MaxOrderRank) != nil { // Look for free cook for current priority
					go s.prepareOrder(s.queue.Remove(idx), ctx)
				}
			}
			s.queueMutex.Unlock()
		}
	}
}

func (s *KitchenSupervisor) getOrderQueue() priorityqueue.IQueue {
	s.queueMutex.Lock()
	defer s.queueMutex.Unlock()
	return s.queue
}

func (s *KitchenSupervisor) prepareOrder(order dto.Order, ctx context.Context) {
	logger.LogMessageF("Preparing order %v", order.OrderID)

	itemChan := make(chan struct{})

	for _, val := range order.Items {
		go s.prepareItem(s.foods[val], itemChan, ctx)
	}

	for len(order.Items) > 0 {
		<-itemChan

		order.Items = order.Items[:len(order.Items)-1] // If we received a signal - pop 1 item
	}

	gateway.Distribute(&order)
}

// prepareItem is a goroutine function which performs infinite search for cook and after this - prepares item
func (s *KitchenSupervisor) prepareItem(item entity.Food, itemChan chan<- struct{}, ctx context.Context) {
	var currCook *cook.Cook		// Holder for cook

	for currCook == nil{
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second): // Wait 1 second before searching for next cook
			currCook = s.findCook(item.Complexity)
		}
	}
	currCook.Prepare(item, itemChan, ctx)	// Immediately start preparing

	logger.LogCookF(currCook.GetID(), "item %d is ready", item.ID)
}

// findCook goes through array of cooks once and return first cook with minimum required rank. nil if no cook found
func (s *KitchenSupervisor) findCook(rankReq int) *cook.Cook {
	s.cooksMutex.Lock()
	defer s.cooksMutex.Unlock()
	for idx, _ := range s.cooks {
		if s.cooks[idx].GetState() == cook.Free && s.cooks[idx].Rank >= rankReq {
			return &s.cooks[idx]
		}
	}
	return nil
}
