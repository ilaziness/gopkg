package main

import (
	"log"

	"github.com/ilaziness/gopkg/crawler/csgo/igxe"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//go c5game.Start()

	igxe.NewIgxe().Fetch()
}
