package main

import (
	"friday/internal/application"
	"log"
)

func main() {
	if err := application.New().Run(); err != nil {
		log.Fatalln(err)
	}
}
