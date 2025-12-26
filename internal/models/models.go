package models


type Event struct {
	ID      string                 `json:"id"`
	Meta    map[string]interface{} `json:"meta"`
	Content []byte                 `json:"content"`
}


