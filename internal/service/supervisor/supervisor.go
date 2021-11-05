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
	queueMutex sync.RWMutex
	queue      priorityqueue.IQueue

	cooksMutex sync.RWMutex
	cooks      []cook.Cook

	foods []entity.Food

	apparatusesMutex sync.RWMutex
	apparatuses      []entity.Apparatus
}

func NewKitchenSupervisor() ISupervisor {
	return &KitchenSupervisor{
		cooks:       cook.NewCooks(repository.GetCooks()),
		foods:       repository.GetFoods(),
		apparatuses: repository.GetApparatuses(),
	}
}

func (s *KitchenSupervisor) AddOrder(order dto.Order) {
	logger.LogSuperF("Got a new order from client #%d", order.ClientID)
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
	queueLookupTick := time.Tick(time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-queueLookupTick:
			s.queueMutex.RLock()
			for idx := 0; idx < s.queue.Len(); idx++ { // Go through all orders in queue
				rankRequired := s.queue.Check(idx).MaxOrderRank // Get rank required for current order
				if s.findCook(rankRequired) {                   // Look for free cook for current rank requirement
					go s.prepareOrder(s.queue.Remove(idx), ctx)
				} else {
					logger.LogSuperF("Currently no free cooks, order from client #%d is hanging", s.queue.Check(idx).ClientID)
				}
			}
			s.queueMutex.RUnlock()
		}
	}
}

func (s *KitchenSupervisor) getOrderQueue() priorityqueue.IQueue {
	s.queueMutex.RLock()
	defer s.queueMutex.RUnlock()
	return s.queue
}

func (s *KitchenSupervisor) prepareOrder(order dto.Order, ctx context.Context) {
	logger.LogMessageF("Preparing order for client #%v", order.ClientID)

	orderChan := make(chan struct{})

	for _, val := range order.Items {
		go s.prepareItem(s.foods[val], orderChan, ctx)
	}

	for len(order.Items) > 0 {
		<-orderChan

		order.Items = order.Items[:len(order.Items)-1] // If we received a signal - pop 1 item
	}

	resp, err := gateway.DistributeV2(&order)
	if err != nil {
		logger.LogMessageF("Could not deliver distribution: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		logger.LogMessage("Delivery was unsuccessful")
	}

	logger.LogMessageF("Suborder for client %d was distributed", order.ClientID)
}

// prepareItem is a goroutine function which performs infinite search for cook and after this - prepares item
func (s *KitchenSupervisor) prepareItem(item entity.Food, itemChan chan<- struct{}, ctx context.Context) {
	var currCook *cook.Cook                          // Holder for cook
	var cookLookupTick = time.Tick(time.Second) // Wait 0.5 seconds before searching for next cook

	for currCook == nil {
		select {
		case <-ctx.Done():
			return
		case <-cookLookupTick:
			s.cooksMutex.RLock()
			for idx := range s.cooks {
				if s.cooks[idx].GetState() == cook.Free && s.cooks[idx].Rank >= item.Complexity {
					if !s.lockApparatus(item) {
						break
					}

					if !s.cooks[idx].Prepare(item, itemChan, ctx) {
						s.unlockApparatus(item)
						continue
					}
					s.unlockApparatus(item)

					s.cooksMutex.RUnlock()
					return
				}
			}
			s.cooksMutex.RUnlock()
		}
	}
}

// findCook goes through array of cooks once and return first cook with minimum required rank. nil if no cook found
func (s *KitchenSupervisor) findCook(rankReq int) bool {
	s.cooksMutex.RLock()
	defer s.cooksMutex.RUnlock()
	for idx := range s.cooks {
		if s.cooks[idx].GetState() == cook.Free && s.cooks[idx].Rank >= rankReq {
			return true
		}
	}
	return false
}

func (s *KitchenSupervisor) lockApparatus(item entity.Food) bool {
	if item.CookingApparatus == "" {
		return true
	}

	for idx, _ := range s.apparatuses {
		if s.apparatuses[idx].Name == item.CookingApparatus {
			if s.apparatuses[idx].IsLocked == true {
				logger.LogWarningF("Tried to allocate already busy apparatus %s for %s", item.CookingApparatus, item.Name)
				return false
			}else{
				s.apparatusesMutex.Lock()
				if s.apparatuses[idx].IsLocked == true {
					logger.LogWarningF("[DOUBLE CHECK] Tried to allocate already busy apparatus %s for %s", item.CookingApparatus, item.Name)
					s.apparatusesMutex.Unlock()

					return false
				}
				s.apparatuses[idx].IsLocked = true
				s.apparatusesMutex.Unlock()
				logger.LogSuperF("%s was locked for %s", item.CookingApparatus, item.Name)
				return true
			}
		}
	}
	logger.LogErrorF("Tried to allocate nonexistent apparatus %s for item", item.CookingApparatus, item.Name)

	return false
}

func (s *KitchenSupervisor) unlockApparatus(item entity.Food) {
	if item.CookingApparatus == "" {
		return
	}

	for idx := range s.apparatuses {
		if s.apparatuses[idx].Name == item.CookingApparatus {
			s.apparatusesMutex.Lock()
			s.apparatuses[idx].IsLocked = false
			s.apparatusesMutex.Unlock()

			logger.LogSuperF("%s was unlocked after %s", item.CookingApparatus, item.Name)
			return
		}
	}
}