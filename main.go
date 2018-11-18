package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)
//import _ "net/http/pprof"

func main() {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	rand.Seed(time.Now().UnixNano())
	instance := new(Instance)
	err := instance.Read("./instances/tai40a")
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	Solve(instance)
	fmt.Println(time.Since(start))
}
