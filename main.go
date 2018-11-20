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
	err := instance.Read("./instances/tai20a")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(calcCost([]int{19, 1, 11, 5, 8, 14, 9, 12, 18, 13, 10, 2, 6, 7, 17, 3, 15, 16, 4, 0}, instance))
	start := time.Now()
	Solve(instance)
	fmt.Println(time.Since(start))
}
