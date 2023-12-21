package main 

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func readFileAndProcess(filename string) ([]*TreeNode, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var trees []*TreeNode

	for _, line := range lines {
		if line == "" {
			continue
		}
		numberStrings := strings.Fields(line)
		var root *TreeNode

		for _, numStr := range numberStrings {
			num, _ := strconv.Atoi(numStr)
			if root == nil {
				root = &TreeNode{Value: num}
			} else {
				root.Insert(num)
			}
		}
		trees = append(trees, root)
	}
	return trees, nil
}