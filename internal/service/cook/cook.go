package cook

import (
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
	"sync"
	"time"
)

type State int

const (
	Free State = iota
	Busy
)

type Cook struct {
	statusMutex sync.Mutex
	status State

	idMutex sync.Mutex
	id int

	currentItemsMutex sync.Mutex
	currentItems int

	entity.Cook
}

func (c *Cook) GetID() int {
	c.idMutex.Lock()
	defer c.idMutex.Unlock()

	return c.id
}

func (c *Cook) GetState() State{
	c.statusMutex.Lock()
	defer c.statusMutex.Unlock()
	return c.status
}

func (c *Cook) SetState(newState State) {
	c.statusMutex.Lock()
	c.status = newState
	c.statusMutex.Unlock()
}

func (c *Cook) Prepare(food entity.Food,foodID int, idChannel chan<- int) {
	c.itemsDecr()
	time.Sleep(time.Second * time.Duration(food.PreparationTime))

	idChannel<-foodID	// send an item prepared signal
	c.itemsIncr()
}

func EntityToService(cookEntities []entity.Cook) []Cook {
	var response []Cook

	for idx, val := range cookEntities {
		response = append(response, Cook{
			status:      Free,
			id:          idx,
			Cook: val,
			currentItems: val.Proficiency,
		})
	}

	return response
}

func (c *Cook) itemsIncr(){
	c.currentItemsMutex.Lock()
	defer c.currentItemsMutex.Unlock()
	c.SetState(Free)
	c.currentItems++
}

func (c *Cook) itemsDecr(){
	c.currentItemsMutex.Lock()
	defer c.currentItemsMutex.Unlock()
	if c.currentItems < 1 {
		logger.LogErrorF("Trying to decrement already zero items!")
	}else if c.currentItems == 1 {
		c.SetState(Busy)
	}
	c.currentItems--
}

