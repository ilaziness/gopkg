package main

import (
	"fmt"
	"log"
)

var V int

func init() {
	log.Println("a plugin init")
}

func GetName() string {
	return "a plugin"
}

func PrintV() {
	log.Println("value of V:", V)
}

var T = Test{c: "c"}

type Test struct {
	A int
	B string
	c string
}

func (t Test) GetC() string {
	return t.c
}

func (t Test) String() string {
	return fmt.Sprintf("%d - %s - %s", t.A, t.B, t.c)
}
