package models

type Room struct {
	Id   int64   `json:"Id"`
	UUID *string `json:"UUID"`
	Name *string `json:"Name"`
}
