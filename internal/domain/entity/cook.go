package entity

type Cook struct {
	Rank int `json:"rank"`
	Proficiency int `json:"proficiency"`
	Name string `json:"name"`
	CatchPhrase string `json:"catch_phrase,omitempty"`
}
