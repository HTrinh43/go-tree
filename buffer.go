package main

import (
	"sync"
)

type WorkItem struct {
	id1 int
	id2 int
}

type Buffer struct {
	mu           sync.Mutex
	condFull     *sync.Cond
	condEmpty    *sync.Cond
	buffer       []WorkItem
	maxSize      int
	closed		 bool
}

func NewBuffer(maxSize int) *Buffer {
	b := &Buffer{
		buffer:  make([]WorkItem, 0, maxSize),
		maxSize: maxSize,
		closed: false,
	}
	b.condFull = sync.NewCond(&b.mu)
	b.condEmpty = sync.NewCond(&b.mu)
	return b
}

func (b *Buffer) Insert(item WorkItem) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.buffer) == b.maxSize{
		b.condFull.Wait()
	}
	b.buffer = append(b.buffer, item)

	b.condEmpty.Signal()
}

func (b *Buffer) Remove() (WorkItem, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.buffer) == 0 && !b.closed{

		b.condEmpty.Wait()
	}
		
	if b.closed{
		b.condEmpty.Broadcast()
		return WorkItem{}, false
	}

	item := b.buffer[0]
	b.buffer = b.buffer[1:]
	b.condFull.Signal()

	return item, true
}


func worker(wg *sync.WaitGroup,trees []*TreeNode, buffer *Buffer, adjMatrix [][]bool) {
	defer wg.Done() 
	for !buffer.closed {
		workItem, ok := buffer.Remove()

		if ok && treesAreEqual(trees[workItem.id1], trees[workItem.id2]) {
			adjMatrix[workItem.id1][workItem.id2] = true
			adjMatrix[workItem.id2][workItem.id1] = true
		}	
	}
}

