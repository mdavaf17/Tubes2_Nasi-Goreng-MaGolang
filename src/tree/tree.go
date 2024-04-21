package tree

import "fmt"

type TreeNode struct {
	Title    string
	URL      string
	Children []*TreeNode
}

func NewTreeNode(title, url string) *TreeNode {
	return &TreeNode{
		Title:    title,
		URL:      url,
		Children: []*TreeNode{},
	}
}

func (n *TreeNode) AddChild(child *TreeNode) {
	n.Children = append(n.Children, child)
}

func PrintTree(node *TreeNode, level int) {
	for i := 0; i < level; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("%s: %s\n", node.Title, node.URL)
	for _, child := range node.Children {
		PrintTree(child, level+1)
	}
}

// func main() {
// root := NewTreeNode("Root", "")

// child1 := NewTreeNode("Child 1", "http://example.com/child1")
// child2 := NewTreeNode("Child 2", "http://example.com/child2")

// child11 := NewTreeNode("Child 1.1", "http://example.com/child1_1")
// child12 := NewTreeNode("Child 1.2", "http://example.com/child1_2")

// child1.AddChild(child11)
// child1.AddChild(child12)

// root.AddChild(child1)
// root.AddChild(child2)

// // Print the tree
// printTree(root, 0)
// }
