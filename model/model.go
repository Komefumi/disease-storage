package model

import "gorm.io/gorm"

type Disease struct {
	gorm.Model
	Name        string `validate:"required,min=1" json:"name"`
	Description string `validate:"required,min=10" json:"description"`
}
