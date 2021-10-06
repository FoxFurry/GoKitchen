package priorityqueue

import (
	"container/heap"
	"github.com/foxfurry/go_kitchen/internal/domain/dto"
	"github.com/foxfurry/go_kitchen/internal/infrastracture/logger"
)

type IQueue interface {
	Less(i, j int) bool
	Swap(i, j int)
	Len() int
	Push(order dto.Order)
	Pop() dto.Order
	Remove(i int) dto.Order
	Check(i int) dto.Order
}

type orderQueue struct {
	pq priorityQueue
}

func NewOrderQueue() IQueue{
	nq := priorityQueue{}
	heap.Init(&nq)

	return &orderQueue{pq: nq}
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

func (o *orderQueue) Remove(i int) dto.Order{
	n := len(o.pq)
	if i >= n {
		logger.LogErrorF("Trying to access order #%d with max len of #%d", n, i)
	}
	return heap.Remove(&o.pq, i).(*Item).Order
}

func (o *orderQueue) Check(i int) dto.Order {
	n := len(o.pq)
	if i >= n {
		logger.LogErrorF("Trying to access order #%d with max len of #%d", n, i)
	}
	return o.pq[n-(i+1)].Order
}