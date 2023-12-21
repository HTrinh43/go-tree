package main

import (
    "sync"
    "fmt"
    "time"
	// "runtime"
)

func findPairs() []WorkItem {
	pairs := make([]WorkItem, 0)
	for _, group := range hashToIndex {
		if len(group) > 1 { // if there's more than one tree with the same hash
			for i := 0; i < len(group); i++ {
				for j := i + 1; j < len(group); j++ {
					pairs = append(pairs, WorkItem{id1: group[i], id2: group[j]})
				}
			}
		}
	}
	return pairs
}

func printGroups(adjMatrix [][]bool) {
	N := len(adjMatrix)
	visited := make([]bool, N)
    count := 0
	for i := 0; i < N; i++ {
		if visited[i] {
			continue
		}

		// For each tree, gather all the equivalent trees
		group := []int{i}
		for j := i + 1; j < N; j++ {
			if adjMatrix[i][j] {
				group = append(group, j)
				visited[j] = true
			}
		}

		// Print the group only if it contains more than one tree
		if len(group) > 1 {
			fmt.Printf("group %d: ", count)
			for _, id := range group {
				fmt.Printf("id%d ", id)
			}
			fmt.Println()
            count += 1
		}
	}
}

func insertWork(buffer *Buffer) {
	for _, group := range hashToIndex {
		if len(group) > 1 { // if there's more than one tree with the same hash
			for i := 0; i < len(group); i++ {
				for j := i + 1; j < len(group); j++ {
					buffer.Insert(WorkItem{id1: group[i], id2: group[j]})
				}
			}
		}
	}
}

func sequentialTreeComparisons(trees []*TreeNode) [][]bool {
    startTime := time.Now() // Capture the start time

	N := len(trees)
    
	// Initialize the 2D adjacency matrix to false
	adjMatrix := make([][]bool, N)
	for i := range adjMatrix {
		adjMatrix[i] = make([]bool, N)
	}

	var wg sync.WaitGroup
	for _, indices := range hashToIndex {
		// If there's only one BST for this hash, set the diagonal element to true
        l := len(indices)
		if l == 1 {
			adjMatrix[indices[0]][indices[0]] = true
			continue
		}
		// For each pair of BSTs with the same hash, compare them
		for i := 0; i < l ; i++ {
			for j := i + 1; j < l ; j++ {
				wg.Add(1)

				go func(indices []int,i, j int) {
					defer wg.Done()
                    if i >= l || j >= l {
                        fmt.Println("Error: Index out of bounds for 'indices' slice.")
                        return
                    }
                    if indices[i] >= len(adjMatrix) || indices[j] >= len(adjMatrix) {
                        fmt.Println("Error: Index out of bounds for 'adjMatrix'.")
                        return
                        }
					if treesAreEqual(trees[indices[i]], trees[indices[j]]) {
						adjMatrix[indices[i]][indices[j]] = true
						adjMatrix[indices[j]][indices[i]] = true
					}
				}(indices, i, j)
			}
		}
	}
	wg.Wait()
    endTime := time.Now() // Capture the end time
    duration := endTime.Sub(startTime).Seconds() // Calculate the duration in seconds
    fmt.Printf("compareTreeTime: %f\n", duration)
	return adjMatrix
}

func insertW(wg *sync.WaitGroup,buffer *Buffer, items []WorkItem, start int){
	// defer wg.Done()
	buffer.Insert(items[start])
}
func parallelTreeComparisons(trees []*TreeNode) [][]bool {
    startTime := time.Now() // Capture the start time
	var wg sync.WaitGroup
	var bwg sync.WaitGroup
	N := len(trees)
    buffer := NewBuffer(compWorkers)

	// Initialize the 2D adjacency matrix to false
	adjMatrix := make([][]bool, N)
	for i := range adjMatrix {
		adjMatrix[i] = make([]bool, N)
	}
	actualWorkers := compWorkers
    pairs := findPairs()

	for j := 0; j < actualWorkers; j++ {
		wg.Add(1)
		go worker(&wg, trees, buffer, adjMatrix)
	}
	for i:=0; i<len(pairs); i++{
		// bwg.Add(1)
		insertW(&bwg,buffer, pairs, i)
	}
	// bwg.Wait()
	buffer.closed=true
	wg.Wait()

    endTime := time.Now() // Capture the end time
    duration := endTime.Sub(startTime).Seconds() // Calculate the duration in seconds
    fmt.Printf("compareTreeTime: %f\n", duration)
	return adjMatrix
}