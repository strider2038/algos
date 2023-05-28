package main

import (
	"log"
)

func main() {
	ui := NewUI()
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
