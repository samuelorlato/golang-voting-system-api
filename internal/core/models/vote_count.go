package models

type VoteCount struct {
	Option int `json:"option"`
	Count  int `json:"count"`
}
