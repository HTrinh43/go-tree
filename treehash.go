package main 

import (
	"sync"
)

func manager(ch chan [2]int, mwg *sync.WaitGroup) {
	defer mwg.Done()
	for result := range ch {
		mapMutex.Lock()
		hashToIndex[result[1]] = append(hashToIndex[result[1]], result[0])
		mapMutex.Unlock()
	}
}

func hashTree(root *TreeNode, index int, ch chan<- [2]int  ) int{
	var sequence []int
	inOrderTraversal(root, &sequence)

	hash := 1
	for _, value := range sequence {
		newValue := value + 2
		hash = (hash*newValue + newValue) % 1000
	}
	if ch != nil {
		ch <- [2]int{index, hash}
	}

	return hash
}

func hashAllTrees(trees []*TreeNode, numWorkers int, ch chan [2]int) {
	var wg sync.WaitGroup
	var	mwg sync.WaitGroup
	if hashWorkers == 1 && dataWorkers == 1 {
		for i, tree := range trees {
        	hash := hashTree(tree, i, nil)
			hashToIndex[hash] = append(hashToIndex[hash], i)
    	}
	} else if hashWorkers > 1 && dataWorkers == 1{
		mwg.Add(1)
		go manager(ch,&mwg)
		for i := 0; i < hashWorkers; i++ {
			wg.Add(1)
			go func(start int, wg *sync.WaitGroup) {
				defer wg.Done()
				for j := start; j < len(trees); j += hashWorkers {
					hashTree(trees[j], j, ch)
				}
			}(i, &wg)
		}
		
		wg.Wait()
		close(ch)	
		mwg.Wait()
	}else if hashWorkers > 1 && dataWorkers == hashWorkers {
		for i := 0; i < hashWorkers; i++ {
			wg.Add(1)
			go func(start int, wg *sync.WaitGroup) {
				defer wg.Done()
				for j := start; j < len(trees); j += hashWorkers {
					hash := hashTree(trees[j], j, nil)
					mapMutex.Lock()
					hashToIndex[hash] = append(hashToIndex[hash], j)
					mapMutex.Unlock()
				}
			}(i, &wg)
		}
		wg.Wait()
		close(ch)	
	} else if hashWorkers > 1 && dataWorkers < hashWorkers && dataWorkers > 1{
		for i := 0; i < hashWorkers; i++ {
			wg.Add(1)
			go func(start int, wg *sync.WaitGroup) {
				defer wg.Done()
				for j := start; j < len(trees); j += hashWorkers {
					hashTree(trees[j], j, ch)
				}
			}(i, &wg)
		}
		for i := 0; i < dataWorkers; i++{
			mwg.Add(1)
			go func(ch chan [2]int, mwg *sync.WaitGroup) {
				defer mwg.Done()
				for result := range ch {
				mapMutex.Lock()
				hashToIndex[result[1]] = append(hashToIndex[result[1]], result[0])
				mapMutex.Unlock()
				}
			}(ch, &mwg)
		}
		wg.Wait()
		close(ch)	
		mwg.Wait()
	}
}

