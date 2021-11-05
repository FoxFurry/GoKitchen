package cook

import (
	"context"
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

type CookSlot struct {
	ItemID   int
	Progress int
	status   State
}

type Cook struct {
	id      int

	itemsMutex   sync.RWMutex
	items        []CookSlot

	currentItemsMutex sync.RWMutex
	currentItems int

	entity.Cook
}

func (c *Cook) GetID() int {
	return c.id
}

func (c *Cook) GetState() State {
	c.currentItemsMutex.RLock()
	defer c.currentItemsMutex.RUnlock()
	if c.currentItems == c.Proficiency {
		return Busy
	}
	return Free
}

func (c *Cook) GetSlots() []CookSlot {
	c.itemsMutex.RLock()
	defer c.itemsMutex.RUnlock()
	return c.items
}

func (c *Cook) Prepare(food entity.Food, isReady chan<- struct{}, ctx context.Context) bool{
	currentFoodID := food.ID
	currentSlot := c.allocateSlot(currentFoodID)

	if currentSlot == -1 {
		logger.LogWarningF("Trying to allocate busy cook. Cook %d | Item %d", c.id, food.ID)
		return false
	}

	logger.LogCookF(c.id, "Preparing item %d in slot %d", currentFoodID, currentSlot)

	var timeTaken float64 = 0
	var progressTick = time.Tick(time.Second)
	var progressComplete = time.After(time.Second * time.Duration(food.PreparationTime))

	for {
		select {
		case <-ctx.Done():
			logger.LogCookF(c.id, "Terminating preparation of %d", currentFoodID)
			return true
		case <-progressTick:
			timeTaken++
			c.items[currentSlot].Progress = int(timeTaken / float64(food.PreparationTime) * 100)

		case <-progressComplete:
			logger.LogCookF(c.id, "Item %d is prepared", currentFoodID)

			c.items[currentSlot] = CookSlot{} // Emptying current slot
			c.decrementItems()

			isReady <- struct{}{}

			return true
		}
	}
}

func (c *Cook) allocateSlot(itemID int) int {
	for idx, val := range c.items {
		if val.status == Free {
			c.itemsMutex.Lock()
			if val.status != Free {	// Double lock check
				c.itemsMutex.Unlock()
				continue
			}else {
				c.items[idx].status = Busy
				c.items[idx].Progress = 0
				c.items[idx].ItemID = itemID

				c.incrementItems()
				c.itemsMutex.Unlock()

				return idx
			}
		}
	}

	return -1
}

func NewCooks(cookEntities []entity.Cook) []Cook {
	var response []Cook

	for idx, val := range cookEntities {
		response = append(response, Cook{
			id:           idx,
			Cook:         val,
			items:        make([]CookSlot, val.Proficiency, val.Proficiency),
			currentItems: 0,
		})
	}

	return response
}

func (c *Cook) incrementItems(){
	c.currentItemsMutex.Lock()
	c.currentItems++
	if c.currentItems > c.Proficiency {
		logger.LogPanic("Current items are more than proficiency")
	}
	c.currentItemsMutex.Unlock()
}

func (c *Cook) decrementItems(){
	c.currentItemsMutex.Lock()
	c.currentItems--
	if c.currentItems < 0 {
		logger.LogPanic("Current items are less than 0")
	}
	c.currentItemsMutex.Unlock()
}