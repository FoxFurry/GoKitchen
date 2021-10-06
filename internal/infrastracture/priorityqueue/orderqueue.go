package priorityqueue

import (
	"container/heap"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
)

type orderQueue struct {
	pq priorityQueue
}

var _ IQueue = &orderQueue{}

func NewOrderQueue() IQueue{
	nq := orderQueue{pq: priorityQueue{}}
	heap.Init(&nq)
	return &nq
}

func (o orderQueue) Less(i, j int) bool {
	return o.pq.Less(i,j)
}

func (o orderQueue) Swap(i, j int) {
	o.pq.Swap(i, j)
}

func (o orderQueue) Len() int {
	return o.pq.Len()
}

func (o *orderQueue) Push(x dto.Order){
	heap.Push(&o.pq, x)
}

func (o *orderQueue) Pop() dto.Order{
	return heap.Pop(&o.pq).(*Item).Order
}
