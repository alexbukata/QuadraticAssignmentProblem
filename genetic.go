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
	for i := 0; i < len(instance.Distances); i++ {
		for j := 0; j < len(instance.Distances); j++ {
			cost += instance.Flows[assignment[i]][assignment[j]] * instance.Distances[i][j]
		}
	}
	return cost
}

func Solve(instance *Instance) Solution {
	populationSize := 100
	population := generateInitialPopulation(instance, populationSize, len(instance.Distances))
	var minIndiv Solution
	min := math.MaxInt64
	for i := 0; i < 200000; i++ {
		parents := selectParents(population, populationSize/2)
		children := generateChildren(parents, populationSize)
		population = append(parents, children...)
		population = mutate(population)

		for _, indiv := range population {
			if indiv.cost < min {
				minIndiv = indiv
				min = indiv.cost
				//fmt.Print(i)
				//fmt.Print(" New best result ")
				//fmt.Println(min)
				//fmt.Println(indiv.String())
			}
		}
	}
	fmt.Println(minIndiv.String())
	return minIndiv
}

func generateInitialPopulation(instance *Instance, populationSize int, geneNumber int) []Solution {
	population := make([]Solution, populationSize)
	instanceInd := 0
	for instanceInd < populationSize{
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
		cost := calcCost(assignment, instance)
		valid := true
		for j := 0; j < instanceInd; j++ {
			if population[j].cost == cost {
				valid = false
				break
			}
		}
		if valid {
			population[instanceInd] = Solution{instance: instance, assignment: assignment, cost: cost}
			instanceInd++
		}
	}
	return population
}

func selectParents(solutions []Solution, number int) []Solution {
	parents := make([]Solution, len(solutions))
	copy(parents, solutions)
	sort.Slice(parents, func(i, j int) bool {
		return parents[i].cost < parents[j].cost
	})
	return parents[:number]
}

func generateChildrenSync(parents []Solution, requiredSize int) []Solution {
	children := make([]Solution, requiredSize*2)
	childrenIndex := 0
	for childrenIndex <= requiredSize {
		firstParentIndex := rand.Intn(len(parents))
		secondParentIndex := rand.Intn(len(parents))
		if firstParentIndex == secondParentIndex {
			continue
		}
		firstParent := parents[firstParentIndex]
		secondParent := parents[secondParentIndex]
		firstChild := crossover(&firstParent, &secondParent)
		valid := true
		for _, parent := range parents {
			if parent.cost == firstChild.cost {
				valid = false
				break
			}
		}
		if valid {
			children[childrenIndex] = firstChild
			childrenIndex++
		}
		secondChild := crossover(&secondParent, &firstParent)
		valid = true
		for _, parent := range parents {
			if parent.cost == secondChild.cost {
				valid = false
				break
			}
		}
		if valid && firstChild.cost != secondChild.cost {
			children[childrenIndex] = secondChild
			childrenIndex++
		}
	}
	return children[:childrenIndex]
}

func generateChildren(parents []Solution, requiredSize int) []Solution {
	children := make([]Solution, requiredSize)
	childrenIndex := 0
	for childrenIndex <= requiredSize {
		childrenChan := make(chan Solution, requiredSize*2)
		var wg sync.WaitGroup
		for i := 0; i < requiredSize; i++ {
			wg.Add(1)
			go doGenerateChildren(parents, childrenChan, &wg)
		}
		wg.Wait()
		close(childrenChan)
		for child := range childrenChan {
			valid := true
			for _, parent := range parents {
				if parent.cost == child.cost {
					valid = false
					break
				}
			}
			for j := 0; j < childrenIndex; j++ {
				if children[j].cost == child.cost {
					valid = false
					break
				}
			}
			if valid {
				children[childrenIndex] = child
				childrenIndex++
				if childrenIndex == requiredSize {
					return children
				}
			}
		}
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
	firstChild := crossover(&firstParent, &secondParent)
	secondChild := crossover(&secondParent, &firstParent)
	childrenChan <- firstChild
	childrenChan <- secondChild
	wg.Done()
}

func crossover(firstParent *Solution, secondParent *Solution) Solution {
	firstParentAssignment := firstParent.assignment
	secondParentAssignment := secondParent.assignment
	newChildGenes := make([]int, len(firstParentAssignment))
	geneIndex := 0
	firstBreak := rand.Intn(len(firstParentAssignment))
	secondBreak := rand.Intn(len(firstParentAssignment))
	for {
		if firstBreak != secondBreak {
			break
		}
		firstBreak = rand.Intn(len(firstParentAssignment))
		secondBreak = rand.Intn(len(firstParentAssignment))
	}
	if firstBreak > secondBreak {
		firstBreak, secondBreak = secondBreak, firstBreak
	}
	for i := 0; i < len(secondParentAssignment); i++ {
		for j := 0; j < firstBreak; j++ {
			if secondParentAssignment[i] == firstParentAssignment[j] {
				newChildGenes[geneIndex] = secondParentAssignment[i]
				geneIndex++
				break
			}
		}
	}
	for i := firstBreak; i <= secondBreak; i++ {
		newChildGenes[geneIndex] = firstParentAssignment[i]
		geneIndex++
	}
	for i := 0; i < len(secondParentAssignment); i++ {
		for j := secondBreak + 1; j < len(firstParentAssignment); j++ {
			if secondParentAssignment[i] == firstParentAssignment[j] {
				newChildGenes[geneIndex] = secondParentAssignment[i]
				geneIndex++
				break
			}
		}
	}
	return Solution{instance: firstParent.instance, assignment: newChildGenes, cost: calcCost(newChildGenes, firstParent.instance)}
}

func mutate(solutions []Solution) []Solution {
	for i := 0; i < len(solutions); i++ {
		if rand.Float32() < 0.1 {
			firstIndex := rand.Intn(len(solutions[i].assignment))
			secondIndex := rand.Intn(len(solutions[i].assignment))
			solutions[i].assignment[firstIndex], solutions[i].assignment[secondIndex] = solutions[i].assignment[secondIndex], solutions[i].assignment[firstIndex]
			solutions[i].cost = calcCost(solutions[i].assignment, solutions[i].instance)
		}
	}
	return solutions
}
