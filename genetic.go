package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

type Solution struct {
	instance   *Instance
	assignment []int
	cost       int
}

func (s *Solution) String() string {
	return fmt.Sprintf("instance:\t%s\nassignment=%v\ncost=%d\n", s.instance.String(), s.assignment, s.cost)
}

func calcCost(assignment []int, instance *Instance) (cost int) {
	for i := 0; i < len(instance.Distances)-1; i++ {
		for j := i + 1; j < len(instance.Distances); j++ {
			cost += instance.Flows[i][j] * instance.Distances[assignment[i]][assignment[j]]
		}
	}
	return cost
}

func Solve(instance *Instance) Solution {
	population := generateInitialPopulation(instance, 2000, len(instance.Distances))
	var minIndiv Solution
	min := math.MaxInt64
	for i := 0; i < 300000; i++ {
		parents := selectParents(population)
		children := generateChildren(parents, 2000)
		population = mutate(children)
		for _, indiv := range population {
			if indiv.cost < min {
				minIndiv = indiv
				min = indiv.cost
			}
		}
	}
	fmt.Println(minIndiv.String())
	return minIndiv
}

func generateInitialPopulation(instance *Instance, populationSize int, geneNumber int) []Solution {
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
		population[i] = Solution{instance: instance, assignment: assignment}
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

func generateChildrenOld(solutions []Solution, requiredSize int) []Solution {
	var children []Solution
	for len(children) <= requiredSize {
		firstParentIndex := rand.Intn(len(solutions))
		secondParentIndex := rand.Intn(len(solutions))
		if firstParentIndex == secondParentIndex {
			continue
		}
		firstParentAssignment := solutions[firstParentIndex].assignment
		secondParentAssignment := solutions[secondParentIndex].assignment
		newChildGenes := make([]int, len(firstParentAssignment))
		copy(newChildGenes, firstParentAssignment[:len(firstParentAssignment)/2])
		geneIndex := len(firstParentAssignment) / 2
		for i := 0; i < len(secondParentAssignment); i++ {
			valid := true
			for j := 0; j < len(newChildGenes); j++ {
				if secondParentAssignment[i] == newChildGenes[j] {
					valid = false
					break
				}
			}
			if valid {
				newChildGenes[geneIndex] = secondParentAssignment[i]
			}
		}
		instance := solutions[firstParentIndex].instance
		newChild := Solution{instance: instance, assignment: newChildGenes, cost: calcCost(newChildGenes, instance)}
		children = append(children, newChild)
	}
	return children
}

func generateChildren(solutions []Solution, requiredSize int) []Solution {
	childrenChan := make(chan Solution, requiredSize*2)
	var children []Solution
	var wg sync.WaitGroup
	for i := 0; i < requiredSize; i++ {
		wg.Add(1)
		go doGenerateChildren(solutions, childrenChan, &wg)
	}
	wg.Wait()
	for i := 0; i < len(childrenChan); i++ {
		children = append(children, <-childrenChan)
	}
	return children
}

func doGenerateChildren(solutions []Solution, childrenChan chan Solution, wg *sync.WaitGroup) {
	firstParentIndex := rand.Intn(len(solutions))
	secondParentIndex := rand.Intn(len(solutions))
	if firstParentIndex == secondParentIndex {
		wg.Done()
		return
	}
	firstParent := solutions[firstParentIndex]
	secondParent := solutions[secondParentIndex]
	instance := firstParent.instance
	firstChild := crossover(&firstParent, &secondParent, instance)
	secondChild := crossover(&secondParent, &firstParent, instance)
	childrenChan <- firstChild
	childrenChan <- secondChild
	wg.Done()
}

func crossover(firstParent *Solution, secondParent *Solution, instance *Instance) Solution {
	firstParentAssignment := firstParent.assignment
	secondParentAssignment := secondParent.assignment
	newChildGenes := make([]int, len(firstParentAssignment))
	copy(newChildGenes, firstParentAssignment[:len(firstParentAssignment)/2])
	geneIndex := len(firstParentAssignment) / 2
	for i := 0; i < len(secondParentAssignment); i++ {
		valid := true
		for j := 0; j < len(newChildGenes); j++ {
			if secondParentAssignment[i] == newChildGenes[j] {
				valid = false
				break
			}
		}
		if valid {
			newChildGenes[geneIndex] = secondParentAssignment[i]
			geneIndex++
		}
	}
	return Solution{instance: instance, assignment: newChildGenes, cost: calcCost(newChildGenes, instance)}
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
