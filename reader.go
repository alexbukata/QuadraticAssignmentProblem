package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instance struct {
	Distances [][]int
	Flows     [][]int
}

func Read(path string) Instance {
	file, e := os.Open(path)
	if e != nil {
		log.Fatal(e)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	numberOfObjectsStr := scanner.Text()
	numberOfObjectsStr = strings.Trim(numberOfObjectsStr, " ")
	numberOfObjects, e := strconv.Atoi(numberOfObjectsStr)
	if e != nil {
		log.Fatal(e)
	}
	distances := readMatrix(numberOfObjects, scanner)
	scanner.Scan()
	flows := readMatrix(numberOfObjects, scanner)
	return Instance{Distances: distances, Flows: flows}
}

func readMatrix(numberOfObjects int, scanner *bufio.Scanner) [][]int {
	matrix := make([][]int, numberOfObjects)
	for i := 0; i < numberOfObjects; i++ {
		scanner.Scan()
		text := scanner.Text()
		distancesStrs := filter(strings.Split(text, " "), func(s string) bool {
			return len(s) > 0
		})
		for j := 0; j < len(distancesStrs); j ++ {
			if j == 0 {
				matrix[i] = make([]int, numberOfObjects)
			}
			matrix[i][j], _ = strconv.Atoi(distancesStrs[j])
		}
	}
	return matrix
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
