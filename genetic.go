package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Solution struct {
	instance   *Instance
	assignment []int
	cost       int
}

func calcCost(assignment []int, instance *Instance) int {
	cost := 0
	for location := 0; location < len(assignment); location++ {
		facility := assignment[location]
		cost += instance.Distances[facility][location] * instance.Flows[facility][location]
	}
	return cost
}

func Solve(instance Instance) []int {
	population := generateInitialPopulation(instance, 1000, len(instance.Distances))
	var minIndiv Solution
	min := 9223372036854775807
	for i := 0; i < 50000; i++ {
		parents := selectParents(population)
		children := generateChildren(parents, 1000)
		population = mutate(children)
		for _, indiv := range population {
			if indiv.cost < min {
				minIndiv = indiv
				min = indiv.cost
			}
		}
	}
	fmt.Printf("%v", minIndiv)
	return []int{}
}

func generateInitialPopulation(instance Instance, populationSize int, geneNumber int) []Solution {
	population := make([]Solution, populationSize)
	for i := 0; i < populationSize; i++ {
		assignment := make([]int, geneNumber)
		//fill
		for j := 0; j < geneNumber; j++ {
			assignment[j] = j
		}
		//shuffle
		for j := 0; j < geneNumber; j++ {
			ix := rand.Intn(geneNumber)
			assignment[j], assignment[ix] = assignment[ix], assignment[j]
		}
		population[i] = Solution{instance: &instance, assignment: assignment}
	}
	return population
}

func selectParents(solutions []Solution) []Solution {
	childs := solutions
	sort.Slice(childs, func(i, j int) bool {
		return childs[i].cost < childs[j].cost
	})
	return childs[len(childs)-int(len(childs)/10):]
}

func generateChildren(solutions []Solution, requiredSize int) []Solution {
	var children []Solution
	for len(children) <= requiredSize {
		firstParentIndex := rand.Intn(len(solutions))
		secondParentIndex := rand.Intn(len(solutions))
		if firstParentIndex == secondParentIndex {
			continue
		}
		firstParentAssignment := solutions[firstParentIndex].assignment
		secondParentAssignment := solutions[secondParentIndex].assignment
		newChildGenes := firstParentAssignment[:len(firstParentAssignment)/2]
	OUTER:
		for i := 0; i < len(secondParentAssignment); i++ {
			for j := 0; j < len(newChildGenes); j++ {
				if secondParentAssignment[i] == newChildGenes[j] {
					continue OUTER
				}
			}
			newChildGenes = append(newChildGenes, secondParentAssignment[i])
		}
		instance := solutions[firstParentIndex].instance
		newChild := Solution{instance: instance, assignment: newChildGenes, cost: calcCost(newChildGenes, instance)}
		children = append(children, newChild)
	}
	return children
}

func mutate(solutions []Solution) []Solution {
	for i := 0; ; i = (i + 1) % len(solutions) {
		if rand.Float32() < 0.01 {
			solutionAssignment := solutions[i].assignment
			firstIndex := rand.Intn(len(solutionAssignment))
			secondIndex := rand.Intn(len(solutionAssignment))
			solutionAssignment[firstIndex], solutionAssignment[secondIndex] = solutionAssignment[secondIndex], solutionAssignment[firstIndex]
			break
		}
	}
	return solutions
}

func populationCost(population []Solution) int {
	sum := 0
	for _, indiv := range population {
		sum += indiv.cost
	}
	return sum
}
