package main

import (
	"log"

	"friday/internal/application"
)

func main() {
	if err := application.New().Run(); err != nil {
		log.Fatalln(err)
	}
}
