package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	instance := new(Instance)
	err := instance.Read("./instances/tai20a")
	if err != nil {
		log.Fatal(err)
	}
	Solve(instance)
}
