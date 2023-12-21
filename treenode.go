package main 

type TreeNode struct {
	Value   int
	Left  *TreeNode
	Right *TreeNode
}

func (node *TreeNode) Insert(value int) {
	if value < node.Value {
		if node.Left == nil {
			node.Left = &TreeNode{Value: value}
		} else {
			node.Left.Insert(value)
		}
	} else if value > node.Value {
		if node.Right == nil {
			node.Right = &TreeNode{Value: value}
		} else {
			node.Right.Insert(value)
		}
	}
}

func inOrderTraversal(root *TreeNode, sequence *[]int) {
	if root == nil {
		return
	}
	inOrderTraversal(root.Left, sequence)
	*sequence = append(*sequence, root.Value)
	inOrderTraversal(root.Right, sequence)
}

func treesAreEqual(root1, root2 *TreeNode) bool {
	if root1 == nil && root2 == nil {
		return true
	}

	if root1 == nil || root2 == nil {
		return false
	}
	return root1.Value == root2.Value &&
		treesAreEqual(root1.Left, root2.Left) &&
		treesAreEqual(root1.Right, root2.Right)
}