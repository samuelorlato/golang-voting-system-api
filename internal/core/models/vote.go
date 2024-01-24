package models

type Vote struct {
	Option string `json:"option" validate:"nonzero"`
}
