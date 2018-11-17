package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var spaces = regexp.MustCompile(`\s+`)

type Instance struct {
	Distances [][]int
	Flows     [][]int
}

func (i *Instance) String() string {
	return fmt.Sprintf("distances:%v\n\tflows:\t%v\n", i.Distances, i.Flows)
}

func (i *Instance) Read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	numberOfObjectsStr := scanner.Text()
	numberOfObjectsStr = strings.Trim(numberOfObjectsStr, " ")
	numberOfObjects, err := strconv.Atoi(numberOfObjectsStr)
	if err != nil {
		return err
	}
	distances, err := readMatrix(numberOfObjects, scanner)
	if err != nil {
		return err
	}
	scanner.Scan()
	flows, err := readMatrix(numberOfObjects, scanner)
	if err != nil {
		return err
	}
	i.Distances = distances
	i.Flows = flows
	return nil
}

func readMatrix(numberOfObjects int, scanner *bufio.Scanner) ([][]int, error) {
	matrix := make([][]int, numberOfObjects)
	for i := 0; i < numberOfObjects; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("scanner failed")
		}
		text := scanner.Text()
		spaceless := strings.TrimSpace(spaces.ReplaceAllString(text, " "))
		distancesStrs := strings.Split(spaceless, " ")
		for _, dist := range distancesStrs {
			number, err := strconv.Atoi(dist)
			if err != nil {
				return nil, err
			}
			matrix[i] = append(matrix[i], number)
		}
	}
	return matrix, nil
}
