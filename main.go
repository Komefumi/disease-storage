package main

import (
	"log"

	"github.com/Komefumi/disease-storage/server"
)

func main() {
	app := server.Setup()

	log.Fatal(app.Listen(":3000"))
}
