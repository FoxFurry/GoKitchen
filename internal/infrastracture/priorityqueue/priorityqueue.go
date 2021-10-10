package priorityqueue

import (
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
)

type item struct {
	dto.Order
	index int
}

type priorityQueue []*item

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq priorityQueue) Len() int { return len(pq) }

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	it := &item{
		Order: x.(dto.Order),
		index: n,
	}
	*pq = append(*pq, it)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil // avoid memory leak
	it.index = -1  // for safety
	*pq = old[0 : n-1]
	return it
}
