package model

type Team struct {
	Data struct {
		Team struct {
			Name    string `json:"name"`
			Players []Player
		} `json:"team"`
	} `json:"data"`
}
