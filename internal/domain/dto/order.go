package dto

import "time"

type Order struct {
	OrderID int   `json:"order_id"`
	TableID int `json:"table_id"`
	WaiterID int `json:"waiter_id"`
	Items   []int `json:"items"`
	Priority int `json:"priority"`
	MaxWait int `json:"max_wait"`
	MaxOrderRank int `json:"-"`
	PickUpTime time.Time `json:"pick_up_time"`
}
