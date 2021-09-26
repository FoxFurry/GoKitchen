package entity

type Food struct {
	ID int `json:"id"`
	Name string `json:"name"`
	PreparationTime int `json:"preparation_time"`
	Complexity int `json:"complexity"`
	CookingApparatus string `json:"cooking_apparatus,omitempty"`
}
