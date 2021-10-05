package priorityqueue

import (
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
)

type Item struct {
	dto.Order
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq *PriorityQueue) Push(x dto.Order) {
	n := len(*pq)
	item := &Item{
		Order: x,
		index: n,
	}
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() dto.Order {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item.Order
}
