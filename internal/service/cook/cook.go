package cook

import (
	"github.com/foxfurry/go_kitchen/internal/domain/entity"
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
	rank int
	proficiency int
	name string
	phrase string
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
	c.SetState(Busy)
	time.Sleep(time.Second * time.Duration(food.PreparationTime))

	idChannel<-foodID
	c.SetState(Free)
}

func EntityToService(cookEntities []entity.Cook) []Cook {
	var response []Cook

	for idx, val := range cookEntities {
		response = append(response, Cook{
			status:      Free,
			id:          idx,
			rank:        val.Rank,
			proficiency: val.Proficiency,
			name:        val.Name,
			phrase:      val.CatchPhrase,
		})
	}

	return response
}

