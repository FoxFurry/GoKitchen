package entity

type Apparatus struct {
	Name string `json:"name"`
	Quantity int `json:"quantity"`
	IsLocked bool `json:"-"`
}
