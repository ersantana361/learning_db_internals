package internal

import (
	"fmt"
)

// BTreeNode represents a node in the B-Tree
type BTreeNode struct {
	ID       string   `json:"id"`
	Keys     []int    `json:"keys"`
	Children []string `json:"children"` // IDs of child nodes
	IsLeaf   bool     `json:"isLeaf"`
	Parent   string   `json:"parent,omitempty"`
}

// BTree represents the B-Tree data structure
type BTree struct {
	Order   int                   `json:"order"` // Maximum number of children
	RootID  string                `json:"rootId"`
	Nodes   map[string]*BTreeNode `json:"nodes"`
	nodeSeq int
}

// NewBTree creates a new B-Tree with the given order
func NewBTree(order int) *BTree {
	if order < 3 {
		order = 3 // Minimum order for a B-Tree
	}
	return &BTree{
		Order:   order,
		RootID:  "",
		Nodes:   make(map[string]*BTreeNode),
		nodeSeq: 0,
	}
}

// generateNodeID creates a unique node ID
func (bt *BTree) generateNodeID() string {
	bt.nodeSeq++
	return fmt.Sprintf("node-%d", bt.nodeSeq)
}

// GetNode returns a node by ID
func (bt *BTree) GetNode(id string) *BTreeNode {
	return bt.Nodes[id]
}

// createNode creates a new node
func (bt *BTree) createNode(isLeaf bool) *BTreeNode {
	node := &BTreeNode{
		ID:       bt.generateNodeID(),
		Keys:     make([]int, 0),
		Children: make([]string, 0),
		IsLeaf:   isLeaf,
	}
	bt.Nodes[node.ID] = node
	return node
}

// Search finds a key in the B-Tree
// Returns (nodeID, keyIndex, found)
func (bt *BTree) Search(key int) (string, int, bool) {
	if bt.RootID == "" {
		return "", -1, false
	}
	return bt.searchNode(bt.RootID, key)
}

func (bt *BTree) searchNode(nodeID string, key int) (string, int, bool) {
	node := bt.Nodes[nodeID]
	if node == nil {
		return "", -1, false
	}

	// Find the first key >= key
	i := 0
	for i < len(node.Keys) && key > node.Keys[i] {
		i++
	}

	// Check if we found the key
	if i < len(node.Keys) && key == node.Keys[i] {
		return nodeID, i, true
	}

	// If leaf, key not found
	if node.IsLeaf {
		return nodeID, i, false
	}

	// Recurse to child
	return bt.searchNode(node.Children[i], key)
}

// Insert adds a key to the B-Tree
func (bt *BTree) Insert(key int) {
	if bt.RootID == "" {
		// Create root node
		root := bt.createNode(true)
		root.Keys = append(root.Keys, key)
		bt.RootID = root.ID
		return
	}

	root := bt.Nodes[bt.RootID]
	if len(root.Keys) == bt.Order-1 {
		// Root is full, need to split
		newRoot := bt.createNode(false)
		newRoot.Children = append(newRoot.Children, bt.RootID)
		root.Parent = newRoot.ID
		bt.RootID = newRoot.ID
		bt.splitChild(newRoot.ID, 0)
		bt.insertNonFull(newRoot.ID, key)
	} else {
		bt.insertNonFull(bt.RootID, key)
	}
}

func (bt *BTree) insertNonFull(nodeID string, key int) {
	node := bt.Nodes[nodeID]

	if node.IsLeaf {
		// Insert key in sorted order
		i := len(node.Keys) - 1
		node.Keys = append(node.Keys, 0)
		for i >= 0 && key < node.Keys[i] {
			node.Keys[i+1] = node.Keys[i]
			i--
		}
		node.Keys[i+1] = key
	} else {
		// Find child to recurse into
		i := len(node.Keys) - 1
		for i >= 0 && key < node.Keys[i] {
			i--
		}
		i++

		child := bt.Nodes[node.Children[i]]
		if len(child.Keys) == bt.Order-1 {
			bt.splitChild(nodeID, i)
			if key > node.Keys[i] {
				i++
			}
		}
		bt.insertNonFull(node.Children[i], key)
	}
}

func (bt *BTree) splitChild(parentID string, childIndex int) {
	parent := bt.Nodes[parentID]
	fullChild := bt.Nodes[parent.Children[childIndex]]

	// Create new node for right half
	newNode := bt.createNode(fullChild.IsLeaf)
	newNode.Parent = parentID

	mid := (bt.Order - 1) / 2

	// Move median key up to parent
	medianKey := fullChild.Keys[mid]

	// Move keys to new node
	newNode.Keys = append(newNode.Keys, fullChild.Keys[mid+1:]...)
	fullChild.Keys = fullChild.Keys[:mid]

	// Move children if not leaf
	if !fullChild.IsLeaf {
		newNode.Children = append(newNode.Children, fullChild.Children[mid+1:]...)
		fullChild.Children = fullChild.Children[:mid+1]

		// Update parent references for moved children
		for _, childID := range newNode.Children {
			bt.Nodes[childID].Parent = newNode.ID
		}
	}

	// Insert median key and new child into parent
	parent.Keys = insertAt(parent.Keys, childIndex, medianKey)
	parent.Children = insertAtStr(parent.Children, childIndex+1, newNode.ID)
}

// Delete removes a key from the B-Tree
func (bt *BTree) Delete(key int) bool {
	if bt.RootID == "" {
		return false
	}

	deleted := bt.deleteFromNode(bt.RootID, key)

	// If root has no keys and has a child, make child the new root
	root := bt.Nodes[bt.RootID]
	if len(root.Keys) == 0 && !root.IsLeaf {
		oldRootID := bt.RootID
		bt.RootID = root.Children[0]
		bt.Nodes[bt.RootID].Parent = ""
		delete(bt.Nodes, oldRootID)
	}

	return deleted
}

func (bt *BTree) deleteFromNode(nodeID string, key int) bool {
	node := bt.Nodes[nodeID]
	minKeys := (bt.Order - 1) / 2

	// Find key position
	i := 0
	for i < len(node.Keys) && key > node.Keys[i] {
		i++
	}

	if node.IsLeaf {
		// Case 1: Key is in leaf node
		if i < len(node.Keys) && node.Keys[i] == key {
			node.Keys = removeAt(node.Keys, i)
			return true
		}
		return false
	}

	if i < len(node.Keys) && node.Keys[i] == key {
		// Case 2: Key is in internal node
		leftChild := bt.Nodes[node.Children[i]]
		rightChild := bt.Nodes[node.Children[i+1]]

		if len(leftChild.Keys) > minKeys {
			// Replace with predecessor
			pred := bt.getPredecessor(node.Children[i])
			node.Keys[i] = pred
			return bt.deleteFromNode(node.Children[i], pred)
		} else if len(rightChild.Keys) > minKeys {
			// Replace with successor
			succ := bt.getSuccessor(node.Children[i+1])
			node.Keys[i] = succ
			return bt.deleteFromNode(node.Children[i+1], succ)
		} else {
			// Merge children
			bt.mergeChildren(nodeID, i)
			return bt.deleteFromNode(node.Children[i], key)
		}
	}

	// Case 3: Key is in subtree
	child := bt.Nodes[node.Children[i]]
	if len(child.Keys) == minKeys {
		bt.fillChild(nodeID, i)
	}

	// After fill, the child index might have changed
	if i > len(node.Keys) {
		return bt.deleteFromNode(node.Children[i-1], key)
	}
	return bt.deleteFromNode(node.Children[i], key)
}

func (bt *BTree) getPredecessor(nodeID string) int {
	node := bt.Nodes[nodeID]
	for !node.IsLeaf {
		node = bt.Nodes[node.Children[len(node.Children)-1]]
	}
	return node.Keys[len(node.Keys)-1]
}

func (bt *BTree) getSuccessor(nodeID string) int {
	node := bt.Nodes[nodeID]
	for !node.IsLeaf {
		node = bt.Nodes[node.Children[0]]
	}
	return node.Keys[0]
}

func (bt *BTree) fillChild(parentID string, childIndex int) {
	parent := bt.Nodes[parentID]
	minKeys := (bt.Order - 1) / 2

	// Try borrowing from left sibling
	if childIndex > 0 {
		leftSibling := bt.Nodes[parent.Children[childIndex-1]]
		if len(leftSibling.Keys) > minKeys {
			bt.borrowFromLeft(parentID, childIndex)
			return
		}
	}

	// Try borrowing from right sibling
	if childIndex < len(parent.Children)-1 {
		rightSibling := bt.Nodes[parent.Children[childIndex+1]]
		if len(rightSibling.Keys) > minKeys {
			bt.borrowFromRight(parentID, childIndex)
			return
		}
	}

	// Merge with sibling
	if childIndex > 0 {
		bt.mergeChildren(parentID, childIndex-1)
	} else {
		bt.mergeChildren(parentID, childIndex)
	}
}

func (bt *BTree) borrowFromLeft(parentID string, childIndex int) {
	parent := bt.Nodes[parentID]
	child := bt.Nodes[parent.Children[childIndex]]
	leftSibling := bt.Nodes[parent.Children[childIndex-1]]

	// Move parent key down to child
	child.Keys = insertAt(child.Keys, 0, parent.Keys[childIndex-1])

	// Move sibling's last key up to parent
	parent.Keys[childIndex-1] = leftSibling.Keys[len(leftSibling.Keys)-1]
	leftSibling.Keys = leftSibling.Keys[:len(leftSibling.Keys)-1]

	// Move child pointer if not leaf
	if !leftSibling.IsLeaf {
		movedChildID := leftSibling.Children[len(leftSibling.Children)-1]
		leftSibling.Children = leftSibling.Children[:len(leftSibling.Children)-1]
		child.Children = insertAtStr(child.Children, 0, movedChildID)
		bt.Nodes[movedChildID].Parent = child.ID
	}
}

func (bt *BTree) borrowFromRight(parentID string, childIndex int) {
	parent := bt.Nodes[parentID]
	child := bt.Nodes[parent.Children[childIndex]]
	rightSibling := bt.Nodes[parent.Children[childIndex+1]]

	// Move parent key down to child
	child.Keys = append(child.Keys, parent.Keys[childIndex])

	// Move sibling's first key up to parent
	parent.Keys[childIndex] = rightSibling.Keys[0]
	rightSibling.Keys = rightSibling.Keys[1:]

	// Move child pointer if not leaf
	if !rightSibling.IsLeaf {
		movedChildID := rightSibling.Children[0]
		rightSibling.Children = rightSibling.Children[1:]
		child.Children = append(child.Children, movedChildID)
		bt.Nodes[movedChildID].Parent = child.ID
	}
}

func (bt *BTree) mergeChildren(parentID string, leftIndex int) {
	parent := bt.Nodes[parentID]
	leftChild := bt.Nodes[parent.Children[leftIndex]]
	rightChild := bt.Nodes[parent.Children[leftIndex+1]]

	// Move parent key down to left child
	leftChild.Keys = append(leftChild.Keys, parent.Keys[leftIndex])

	// Move all keys from right child to left child
	leftChild.Keys = append(leftChild.Keys, rightChild.Keys...)

	// Move all children from right child to left child
	if !leftChild.IsLeaf {
		for _, childID := range rightChild.Children {
			bt.Nodes[childID].Parent = leftChild.ID
		}
		leftChild.Children = append(leftChild.Children, rightChild.Children...)
	}

	// Remove key and child pointer from parent
	parent.Keys = removeAt(parent.Keys, leftIndex)
	parent.Children = removeAtStr(parent.Children, leftIndex+1)

	// Delete right child
	delete(bt.Nodes, rightChild.ID)
}

// RangeSearch finds all keys in the range [start, end]
func (bt *BTree) RangeSearch(start, end int) []int {
	result := []int{}
	if bt.RootID == "" {
		return result
	}
	bt.rangeSearchNode(bt.RootID, start, end, &result)
	return result
}

func (bt *BTree) rangeSearchNode(nodeID string, start, end int, result *[]int) {
	node := bt.Nodes[nodeID]

	i := 0
	for i < len(node.Keys) && node.Keys[i] < start {
		i++
	}

	for i < len(node.Keys) && node.Keys[i] <= end {
		if !node.IsLeaf {
			bt.rangeSearchNode(node.Children[i], start, end, result)
		}
		*result = append(*result, node.Keys[i])
		i++
	}

	if !node.IsLeaf && i < len(node.Children) {
		bt.rangeSearchNode(node.Children[i], start, end, result)
	}
}

// Clone creates a deep copy of the B-Tree
func (bt *BTree) Clone() *BTree {
	clone := &BTree{
		Order:   bt.Order,
		RootID:  bt.RootID,
		Nodes:   make(map[string]*BTreeNode),
		nodeSeq: bt.nodeSeq,
	}
	for id, node := range bt.Nodes {
		clone.Nodes[id] = &BTreeNode{
			ID:       node.ID,
			Keys:     append([]int{}, node.Keys...),
			Children: append([]string{}, node.Children...),
			IsLeaf:   node.IsLeaf,
			Parent:   node.Parent,
		}
	}
	return clone
}

// Helper functions
func insertAt(slice []int, index int, value int) []int {
	slice = append(slice, 0)
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

func insertAtStr(slice []string, index int, value string) []string {
	slice = append(slice, "")
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

func removeAt(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func removeAtStr(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}
