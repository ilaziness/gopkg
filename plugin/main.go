package main

import (
	"fmt"
	"log"
	"plugin"
)

// go build -buildmode=plugin -buildvcs=false -o a_plugin.so ./a/

type TestI interface {
	fmt.Stringer
	GetC() string
}

func main() {
	pl, err := plugin.Open("a_plugin.so")
	if err != nil {
		log.Fatalln(err)
	}

	v, err := pl.Lookup("V")
	if err != nil {
		log.Fatalln(err)
	}

	// *v.(*int) = 10
	nv := v.(*int)
	*nv = 10

	name1, err := pl.Lookup("PrintV")
	if err != nil {
		log.Fatalln(err)
	}
	name1.(func())()
	log.Println(*nv)

	name, err := pl.Lookup("GetName")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("plugin name: %s\n", name.(func() string)())

	newTest, err := pl.Lookup("T")
	if err != nil {
		log.Fatalln(err)
	}
	t := newTest.(TestI)
	log.Printf("%s\n", t)
	log.Printf("c val: %s\n", t.GetC())
}
