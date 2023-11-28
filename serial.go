package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const USAGE = "serial [inputFile.txt] [outputFile] [board_size_X] [board_size_Y] [num_processes]"

func main() {

	/** Validating command line arguments **/

	var default_procs int = 1
	// check command line arguments
	var procs *int = &default_procs
	// default number of processes is 1
	var numArgs int = len(os.Args)
	if numArgs < 5 || numArgs > 6 {
		fmt.Println(USAGE)
		return
	}
	// process the number of processes
	if numArgs == 6 {
		numprocs, procerror := strconv.Atoi(os.Args[5])
		if procerror != nil {
			fmt.Println("did not pass valid num_process argument: " + os.Args[5])
		} else {
			*procs = numprocs
		}
	}
	// process the rest of the command line arguments
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	xDim, xErr := strconv.Atoi(os.Args[3])
	yDim, yErr := strconv.Atoi(os.Args[4])

	if xErr != nil || yErr != nil {
		fmt.Println("invalid x and y board dimensions: x= " + os.Args[3] + "  |  y=" + os.Args[4])
	}

	/** **/

	// keep the previous state of the board
	var prev_life [][]int = make([][]int, xDim+2)
	var life [][]int = make([][]int, yDim)
	for i := 0; i < xDim; i++ {
		prev_life[i] = make([]int, yDim+2)
		life[i] = make([]int, yDim)
	}
	// initialize last 2 rows of prev_life
	prev_life[xDim] = make([]int, yDim+2)
	prev_life[xDim+1] = make([]int, yDim+2)

	readBoard(&life, xDim, yDim, inputFile)
	fmt.Printf("value at %d,%d is %d\n", 27, 16, life[27][16])
	fmt.Printf("input=%s, output=%s, X=%d, Y=%d, procs=%d\n", inputFile, outputFile, xDim, yDim, *procs)
}

/*
reads a file where each line is a 0-indexed coordinate on the board separated by a comma,
i.e.

	23,42
	1,4
	2,65
	...
*/
func readBoard(board *[][]int, xDim int, yDim int, path string) {

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
		// validate coordinate dimensions. Coordinates are 0-indexed
		if x < 0 || x > xDim-1 || y < 0 || y > yDim-1 {
			fmt.Printf(`invalid board coordinate: %d,%d
Must be 0-indexed and within the bounds(exclusive): %dx%d
			`, x, y, xDim, yDim)
			fmt.Println(xErr.Error())
			fmt.Println(yErr.Error())
			os.Exit(1)
		}

		// copy the coordinate into the life array
		(*board)[x][y] = 1
		lineNum++
	}
}
