package priorityqueue

import (
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
)

type IQueue interface {
	Less(i, j int) bool
	Swap(i, j int)
	Len() int
	Push(x dto.Order)
	Pop() dto.Order
}

type Item struct {
	dto.Order
	index int
}

type priorityQueue []*Item

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq priorityQueue) Len() int { return len(pq) }

func (pq *priorityQueue) Push(x dto.Order) {
	n := len(*pq)
	item := &Item{
		Order: x,
		index: n,
	}
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() dto.Order {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item.Order
}
