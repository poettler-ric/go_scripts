package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	outFile = flag.String("out", "", "output file")
)

func readLatinSquare(file string) (square [][]int, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed opening %v: %v", file, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed reading csv: %v", err)
	}

	square = make([][]int, len(data))
	for r, row := range data {
		square[r] = make([]int, len(row))
		for c, cell := range row {
			if cell == "" {
				square[r][c] = 0
			} else {
				square[r][c], err = strconv.Atoi(cell)
				if err != nil {
					return nil, fmt.Errorf(
						"failed converting to \"%v\" int: %v",
						cell,
						err)
				}
			}
		}
	}
	return
}

func writeLatinSquare(square [][]int, file string) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed creating %v: %v", file, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, row := range square {
		record := make([]string, len(row))
		for i := 0; i < len(row); i++ {
			record[i] = strconv.Itoa(row[i])
		}
		w.Write(record)
	}
	w.Flush()
	return
}

func validateSquare(square [][]int) (result bool) {
	// check dimension
	dimension := len(square)
	for _, row := range square {
		if len(row) != dimension {
			return false
		}
	}

	// check rows and columns only contain unique elements
	for i := 0; i < dimension; i++ {
		// check row
		seen := make(map[int]bool)
		nonZeros := 0
		for j := 0; j < dimension; j++ {
			if square[i][j] != 0 {
				seen[square[i][j]] = true
				nonZeros++
			}
		}
		if len(seen) != nonZeros {
			return false
		}

		// check column
		seen = make(map[int]bool)
		nonZeros = 0
		for j := 0; j < dimension; j++ {
			if square[j][i] != 0 {
				seen[square[j][i]] = true
				nonZeros++
			}
		}
		if len(seen) != nonZeros {
			return false
		}
	}

	// numbers must be between 0 and dimension (inclusive)
	for _, row := range square {
		for _, cell := range row {
			if cell < 0 || cell > dimension {
				return false
			}
		}
	}
	return true
}

func getPossibleElements(square [][]int, row, column int) (result []int) {
	dimension := len(square)

	seen := make(map[int]bool)
	for i := 0; i < dimension; i++ {
		seen[square[row][i]] = true
		seen[square[i][column]] = true
	}

	result = make([]int, 0, dimension-len(seen)+1)
	for i := 1; i <= dimension; i++ {
		if _, ok := seen[i]; !ok {
			result = append(result, i)
		}
	}
	return
}

func copySquare(square [][]int) (result [][]int) {
	result = make([][]int, len(square))
	for i, row := range square {
		rowCopy := make([]int, len(row))
		copy(rowCopy, row)
		result[i] = rowCopy
	}
	return
}

func nextField(dimension, row, column int) (newRow, newColumn int, err error) {
	newRow = row
	newColumn = column + 1
	if newColumn >= dimension {
		newRow++
		newColumn = newColumn % dimension
	}
	if newRow >= dimension {
		return -1, -1, errors.New("reached end")
	}
	return
}

func solveSquare(square [][]int, row, column int) (result [][]int, ok bool) {
	dimension := len(square)
	newRow, newColumn := row, column
	// determine next field to solve
	for {
		var err error
		newRow, newColumn, err = nextField(dimension, newRow, newColumn)
		// there are no fields left ->  we found a solution
		if err != nil {
			return square, true
		}
		// we found an empty field -> try to solve it
		if square[newRow][newColumn] == 0 {
			break
		}
	}
	// try all possibilities for the field
	for _, i := range getPossibleElements(square, newRow, newColumn) {
		squareCopy := copySquare(square)
		squareCopy[newRow][newColumn] = i
		// try to solve the square with the currecnt possibility
		solution, ok := solveSquare(squareCopy, newRow, newColumn)
		if ok {
			return solution, true
		}
	}
	return nil, false
}

func printSquare(square [][]int) {
	for _, row := range square {
		for _, cell := range row {
			fmt.Printf("%3d", cell)
		}
		fmt.Print("\n")
	}
}

func main() {
	flag.Parse()

	square, err := readLatinSquare(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	solution, ok := solveSquare(square, 0, 0)
	if !ok {
		fmt.Println("couldn't find a solution")
	} else {
		printSquare(solution)
		if *outFile != "" {
			writeLatinSquare(solution, *outFile)
		}
	}
}
