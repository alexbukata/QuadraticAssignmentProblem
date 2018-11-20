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
		//fmt.Print("Generation ")
		//fmt.Println(i)
		parents := selectParents(population, populationSize/2)
		children := generateChildren(parents, populationSize)
		population = append(parents, children...)
		population = mutate(population)

		//repeats := 0
		//for j := 0; j < len(population)-1; j++ {
		//	for k := j + 1; k < len(population); k++ {
		//		first := population[j]
		//		second := population[k]
		//		if first.cost == second.cost {
		//			repeats++
		//		}
		//	}
		//}
		//fmt.Print("Repeats ")
		//fmt.Println(repeats)

		for _, indiv := range population {
			if indiv.cost < min {
				minIndiv = indiv
				min = indiv.cost
				fmt.Print(i)
				fmt.Print(" New best result ")
				fmt.Println(min)
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
		cost := calcCost(assignment, instance)
		for j := 0; j < i; j++ {
			if population[i].cost == cost {
				i--
				continue
			}
		}
		population[i] = Solution{instance: instance, assignment: assignment, cost: cost}
	}
	return population
}

func selectParents(solutions []Solution, number int) []Solution {
	parents := make([]Solution, len(solutions))
	copy(parents, solutions)
	sort.Slice(parents, func(i, j int) bool {
		return parents[i].cost < parents[j].cost
	})
	return parents[len(parents)-number:]
}

func generateChildrenOld(parents []Solution, requiredSize int) []Solution {
	children := make([]Solution, requiredSize)
	childrenIndex := 0
	for childrenIndex <= requiredSize {
		childrenChan := make(chan Solution, 20)
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
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
	instance := firstParent.instance
	firstChild := crossover(&firstParent, &secondParent, instance)
	secondChild := crossover(&secondParent, &firstParent, instance)
	childrenChan <- firstChild
	childrenChan <- secondChild
	wg.Done()
}

func crossoverOld(firstParent *Solution, secondParent *Solution, instance *Instance) Solution {
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

func crossover(firstParent *Solution, secondParent *Solution, instance *Instance) Solution {
	firstParentAssignment := firstParent.assignment
	secondParentAssignment := secondParent.assignment
	newChildGenes := make([]int, len(firstParentAssignment))
	geneIndex := 0
	firstBreak := rand.Intn(len(firstParentAssignment)/2-1) + 1
	secondBreak := rand.Intn(len(firstParentAssignment)/2-2) + len(firstParentAssignment)/2 + 1
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
	return Solution{instance: instance, assignment: newChildGenes, cost: calcCost(newChildGenes, instance)}
}

func mutate(solutions []Solution) []Solution {
	for i := 0; i < len(solutions); i++ {
		if rand.Float32() < 0.1 {
			solutionAssignment := solutions[i].assignment
			firstIndex := rand.Intn(len(solutionAssignment))
			secondIndex := rand.Intn(len(solutionAssignment))
			solutionAssignment[firstIndex], solutionAssignment[secondIndex] = solutionAssignment[secondIndex], solutionAssignment[firstIndex]
			solutions[i].cost = calcCost(solutionAssignment, solutions[i].instance)
		}
	}
	return solutions
}
