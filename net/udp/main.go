package main

import "log"

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	simple()

	sameCS()

	multicast()

	boradcast()
}
