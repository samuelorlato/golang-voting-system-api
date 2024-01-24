package models

type VoteCount struct {
	Count               int     `json:"count"`
	PercentageFromTotal float64 `json:"percentageFromTotal"`
}
