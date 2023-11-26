package main

import (
	"fmt"
	"os"
	"strconv"
)

const USAGE = "serial [input_file.txt] [output_file] [board_size_X] [board_size_Y] [num_processes]"

func main() {
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
	input_file := os.Args[1]
	output_file := os.Args[2]
	x_dim, xerror := strconv.Atoi(os.Args[3])
	y_dim, yerror := strconv.Atoi(os.Args[4])

	if xerror != nil || yerror != nil {
		fmt.Println("invalid x and y board dimensions: x= " + os.Args[3] + "  |  y=" + os.Args[4])
	}

	fmt.Printf("input=%s, output=%s, X=%d, Y=%d, procs=%d\n", input_file, output_file, x_dim, y_dim, *procs)
}
