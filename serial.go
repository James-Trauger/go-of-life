package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const USAGE = "serial [inputFile.txt] [outputFile] [board_size_X] [board_size_Y] [generations] [num_processes]"

func main() {

	/** Validating command line arguments **/

	var default_procs int = 1
	// default number of processes is 1
	var numArgs int = len(os.Args)
	if numArgs < 6 || numArgs > 7 {
		fmt.Println(USAGE)
		return
	}

	procs := func() int {
		if numArgs == 7 {
			numprocs, procerr := strconv.Atoi(os.Args[6])
			if procerr != nil {
				fmt.Printf("\"%s\" is not pass valid number of processors\n", os.Args[6])
				os.Exit(2)
			}
			return numprocs
		} else {
			return default_procs
		}
	}()

	// process the rest of the command line arguments
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	xdim, xErr := strconv.Atoi(os.Args[3])
	ydim, yErr := strconv.Atoi(os.Args[4])
	if xErr != nil || yErr != nil {
		fmt.Println("invalid x and y board ydimensions: x= " + os.Args[3] + "  |  y=" + os.Args[4])
	}

	generations, generr := strconv.Atoi(os.Args[5])
	if generr != nil {
		fmt.Printf("\"%s\" is not a valid number of generations\n", os.Args[5])
	}

	// keep the previous state of the board
	var prev_life [][]int = make([][]int, xdim+2)
	var life [][]int = make([][]int, ydim)
	// initialize every row
	for i := 0; i < xdim; i++ {
		prev_life[i] = make([]int, ydim+2)
		life[i] = make([]int, ydim)
	}
	// initialize last 2 rows of prev_life
	prev_life[xdim] = make([]int, ydim+2)
	prev_life[xdim+1] = make([]int, ydim+2)

	readBoard(&life, xdim, ydim, inputFile)

	// run throught the board for the specified number of generations
	//for := range()

	writeBoard(&life, xdim, ydim, outputFile)
	//fmt.Printf("value at %d,%d is %d\n", 27, 16, life[27][16])
	fmt.Printf("input=%s, output=%s, X=%d, Y=%d, generations=%d, procs=%d\n", inputFile, outputFile, xdim, ydim, generations, procs)
}

/*
reads a file where each line is a 0-indexed coordinate on the board separated by a comma,
i.e.

	23,42
	1,4
	2,65
	...
*/
func readBoard(board *[][]int, xdim int, ydim int, path string) {

	input, err := os.Open(path)
	defer input.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	// each line is formatted as such: 32,54
	// where 32 and 54 are arbitrary row and column numbers
	r := regexp.MustCompile(`(\d+),(\d+)`)
	// for reading the file line by line
	scanner := bufio.NewScanner(input)

	// read each point on the board into the board array
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		match := r.FindStringSubmatch(line)
		x, xErr := strconv.Atoi(match[1])
		y, yErr := strconv.Atoi(match[2])
		// validate the coordinates are integers
		if xErr != nil || yErr != nil {
			fmt.Printf("invalid line format on line %d. Must be in the format: x,y where x and y are integers\nfound: %s\n", lineNum, line)
			os.Exit(1)
		}
		// validate coordinate ydimensions. Coordinates are 0-indexed
		if x < 0 || x > xdim-1 || y < 0 || y > ydim-1 {
			fmt.Printf(`invalid board coordinate: %d,%d
Must be 0-indexed and within the bounds(exclusive): %dx%d
			`, x, y, xdim, ydim)
			fmt.Println(xErr.Error())
			fmt.Println(yErr.Error())
			os.Exit(1)
		}

		// copy the coordinate into the life array
		(*board)[x][y] = 1
		lineNum++
	}
}

func writeBoard(board *[][]int, xdim int, ydim int, output_file string) {
	// create the file
	output, err := os.Create(output_file + ".csv")
	defer output.Close()
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	// write the file in row major order
	// iterate over each row and print only 1s
	for i := 0; i < xdim; i++ {
		for j := 0; j < ydim; j++ {
			if (*board)[i][j] == 1 {
				// write the record
				fmt.Fprintf(output, "%d,%d\n", i, j)
			}
		}
	}
}
