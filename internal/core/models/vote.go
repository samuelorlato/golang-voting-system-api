package models

type Vote struct {
	Option int `json:"option" validate:"nonzero"`
}