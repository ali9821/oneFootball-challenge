package model

type PlayerData struct {
	Id    string   `json:"id"`
	Name  string   `json:"name"`
	Age   string   `json:"age"`
	Teams []string `json:"teams"`
}
