package model

import "gorm.io/gorm"

// dbOpened.Create(&Disease{Name: "ProtoType Disease", Description: "Non real disease, made as a model to perform operations with"})

var PrototypeDisease = Disease{Name: "ProtoType Disease", Description: "Non real disease, made as a model to perform operations with"}

type Disease struct {
	gorm.Model
	Name        string `validate:"required,min=1" json:"name"`
	Description string `validate:"required,min=10" json:"description"`
}
