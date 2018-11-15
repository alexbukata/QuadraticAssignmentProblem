package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	instance := Read("./instances/tai20a")
	Solve(instance)
}
