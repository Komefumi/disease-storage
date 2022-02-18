package database

import (
	"fmt"
	"os"

	"github.com/Komefumi/disease-storage/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	err := os.Remove("test.db")
	if err != nil {
		fmt.Println(err)
	}
	dbOpened, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbOpened.AutoMigrate(&model.Disease{})

	// dbOpened.Create(&Disease{Name: "ProtoType Disease", Description: "Non real disease, made as a model to perform operations with"})

	DB = dbOpened
}
