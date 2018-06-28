package models

type Room struct {
	Id   int64   `json:"key"`
	UUID *string `json:"value"`
	Name *string `json:"text"`
}
