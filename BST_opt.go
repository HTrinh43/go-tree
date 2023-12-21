package main 

import (
	"flag"
	"fmt"
    "os"
	"sync"
	"time"
)

var (
	hashWorkers  int
	dataWorkers  int
	compWorkers  int
	isPrint  int
	trees 		[]*TreeNode
	hashToIndex  = make(map[int][]int)
	mapMutex     = &sync.RWMutex{}
	hashMutexMap = make(map[int]*sync.Mutex)
	
)




func main() {
	// Define command-line flags.
	var inputFile string

	flag.IntVar(&hashWorkers, "hash-workers", 1, "Number of threads for hashing")
	flag.IntVar(&dataWorkers, "data-workers", 1, "Number of threads for data processing")
	flag.IntVar(&compWorkers, "comp-workers", 1, "Number of threads for computations")
	flag.IntVar(&isPrint, "print", 1, "Print results")
	flag.StringVar(&inputFile, "input", "", "Path to an input file")

	// Parse the flags.
	flag.Parse()

	if inputFile == "" {
		fmt.Println("Please provide an input file using the -input parameter.")
		os.Exit(1)
	}

	trees, err := readFileAndProcess(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ch := make(chan [2]int, len(trees))

    startTime := time.Now() // Capture the start time
	// Hashing
    hashAllTrees(trees, hashWorkers, ch)

    endTime := time.Now() // Capture the end time
    duration := endTime.Sub(startTime).Seconds() // Calculate the duration in seconds
    fmt.Printf("hashGroupTime: %f\n", duration)
	if isPrint == 1{
		printHash()
	}
	var adj [][]bool
	if compWorkers == 1 {
		adj = sequentialTreeComparisons(trees)
	} else if compWorkers > 1 {
		adj = parallelTreeComparisons(trees)
	}
	if isPrint == 1{
		printGroups(adj)
	}
}

func printHash(){
	for hash, indices := range hashToIndex {
		fmt.Printf("hash%d: ", hash)
		for _, index := range indices {
			fmt.Printf("%d ", index)
		}
		fmt.Println()
	}	
	fmt.Println()
}