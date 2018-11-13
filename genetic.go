package main

import (
	"fmt"
	"math/rand"
)

func Solve(instance Instance) []int {
	population := generateInitialPopulation(1000, len(instance.Distances))
	fmt.Printf("%v", population)
	return []int{}
}

func generateInitialPopulation(populationSize int, geneNumber int) [][]int {
	population := make([][]int, populationSize)
	for i := 0; i < populationSize; i++ {
		population[i] = make([]int, geneNumber)
		//fill
		for j := 0; j < geneNumber; j++ {
			population[i][j] = j
		}
		//shuffle
		for j := 0; j < geneNumber; j++ {
			ix := rand.Intn(geneNumber)
			population[i][j], population[i][ix] = population[i][ix], population[i][j]
		}
	}
	return population
}
