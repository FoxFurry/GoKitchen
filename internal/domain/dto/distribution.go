package dto

type cookingDetail struct {
	FoodID int `json:"food_id"`
	CookID int `json:"cook_id"`
}

type Distribution struct {
	Order

	CookingTime int                `json:"cooking_time"`
	CookingDetails []cookingDetail `json:"cooking_details"`
}
