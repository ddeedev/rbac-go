package main

import (
	"log"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}
