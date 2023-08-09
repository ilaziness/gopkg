package main

import (
	"github.com/ilaziness/gopkg/crawler/csgo/c5game"
	"github.com/ilaziness/gopkg/crawler/csgo/igxe"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go c5game.Start()

	igxe.NewIgxe().Fetch()
}
