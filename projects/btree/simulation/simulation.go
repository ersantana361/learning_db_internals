package simulation

import (
	"fmt"

	"github.com/ersantana/db-internals/packages/protocol"
	"github.com/ersantana/db-internals/packages/simulation/engine"
	"github.com/ersantana/db-internals/projects/btree/internal"
)

// BTreeSimulation implements the simulation.Simulation interface
type BTreeSimulation struct {
	tree        *internal.BTree
	steps       []engine.Step
	currentStep int
	operation   string
	operand     int
	searchPath  []string
}

// NewBTreeSimulation creates a new B-Tree simulation
func NewBTreeSimulation() *BTreeSimulation {
	return &BTreeSimulation{
		tree:        internal.NewBTree(4), // Order 4 B-Tree (max 3 keys per node)
		steps:       make([]engine.Step, 0),
		currentStep: -1,
	}
}

// Name returns the simulation name
func (sim *BTreeSimulation) Name() string {
	return "B-Tree"
}

// Description returns the simulation description
func (sim *BTreeSimulation) Description() string {
	return "Interactive visualization of B-Tree operations including insert, search, and delete"
}

// Initialize sets up the simulation with given config
func (sim *BTreeSimulation) Initialize(config map[string]interface{}) error {
	// Reset tree
	order := 4
	if o, ok := config["order"].(float64); ok {
		order = int(o)
	}
	sim.tree = internal.NewBTree(order)

	// Pre-populate if specified
	if keys, ok := config["initialKeys"].([]interface{}); ok {
		for _, k := range keys {
			if key, ok := k.(float64); ok {
				sim.tree.Insert(int(key))
			}
		}
	}

	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	return nil
}

// Reset returns the simulation to initial state
func (sim *BTreeSimulation) Reset() error {
	sim.tree = internal.NewBTree(sim.tree.Order)
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.searchPath = nil
	return nil
}

// GenerateSteps returns all steps for current simulation
func (sim *BTreeSimulation) GenerateSteps() []engine.Step {
	return sim.steps
}

// CurrentStep returns current step index
func (sim *BTreeSimulation) CurrentStep() int {
	return sim.currentStep
}

// ExecuteStep executes a specific step
func (sim *BTreeSimulation) ExecuteStep(index int) engine.StepResult {
	if index < 0 || index >= len(sim.steps) {
		return engine.StepResult{
			Success: false,
			Error:   engine.ErrInvalidStepIndex,
		}
	}

	sim.currentStep = index
	step := sim.steps[index]

	return engine.StepResult{
		Success:     true,
		Highlights:  step.Highlights,
		Data:        sim.GetVisualizationData(),
		Description: step.Description,
	}
}

// CanStepForward returns true if can advance
func (sim *BTreeSimulation) CanStepForward() bool {
	return sim.currentStep < len(sim.steps)-1
}

// CanStepBackward returns true if can go back
func (sim *BTreeSimulation) CanStepBackward() bool {
	return sim.currentStep > 0
}

// GetState returns current tree state
func (sim *BTreeSimulation) GetState() interface{} {
	return map[string]interface{}{
		"order":       sim.tree.Order,
		"rootId":      sim.tree.RootID,
		"nodes":       sim.tree.Nodes,
		"operation":   sim.operation,
		"currentStep": sim.currentStep,
		"totalSteps":  len(sim.steps),
	}
}

// GetVisualizationData returns data for rendering
func (sim *BTreeSimulation) GetVisualizationData() map[string]interface{} {
	nodes := make(map[string]interface{})
	for id, node := range sim.tree.Nodes {
		nodes[id] = map[string]interface{}{
			"id":       node.ID,
			"keys":     node.Keys,
			"children": node.Children,
			"isLeaf":   node.IsLeaf,
			"parent":   node.Parent,
		}
	}

	return map[string]interface{}{
		"nodes":  nodes,
		"rootId": sim.tree.RootID,
		"order":  sim.tree.Order,
		"path":   sim.searchPath,
	}
}

// PrepareInsert generates steps for an insert operation
func (sim *BTreeSimulation) PrepareInsert(key int) {
	sim.operation = "insert"
	sim.operand = key
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.searchPath = nil

	// Clone tree for step-by-step visualization
	treeCopy := sim.tree.Clone()

	// Step 1: Start
	sim.addStep(
		fmt.Sprintf("Insert %d", key),
		fmt.Sprintf("Starting insertion of key %d into the B-Tree", key),
		[]protocol.Highlight{},
		sim.cloneTreeData(treeCopy),
	)

	if treeCopy.RootID == "" {
		// Empty tree
		sim.addStep(
			"Create Root",
			fmt.Sprintf("Tree is empty. Creating root node with key %d", key),
			[]protocol.Highlight{{Type: "node", ID: "new-root", Color: "#10b981", Animation: "pulse"}},
			sim.cloneTreeData(treeCopy),
		)
		treeCopy.Insert(key)
		sim.addStep(
			"Insertion Complete",
			fmt.Sprintf("Key %d inserted as root", key),
			[]protocol.Highlight{{Type: "node", ID: treeCopy.RootID, Color: "#10b981", Animation: "pulse"}},
			sim.cloneTreeData(treeCopy),
		)
	} else {
		// Find insertion path
		path := sim.findInsertionPath(treeCopy, key)

		// Traverse to insertion point
		for i, nodeID := range path {
			node := treeCopy.Nodes[nodeID]
			sim.addStep(
				fmt.Sprintf("Traverse to %s", nodeID),
				fmt.Sprintf("Examining node with keys %v. Looking for position to insert %d", node.Keys, key),
				[]protocol.Highlight{{Type: "node", ID: nodeID, Color: "#3b82f6", Animation: "pulse"}},
				sim.cloneTreeData(treeCopy),
			)

			if i == len(path)-1 && node.IsLeaf {
				if len(node.Keys) < treeCopy.Order-1 {
					// Simple insert
					sim.addStep(
						"Insert Key",
						fmt.Sprintf("Leaf node has space. Inserting key %d", key),
						[]protocol.Highlight{
							{Type: "node", ID: nodeID, Color: "#10b981", Animation: "pulse"},
							{Type: "key", ID: fmt.Sprintf("%d", key), Color: "#10b981", Animation: "pulse"},
						},
						sim.cloneTreeData(treeCopy),
					)
				} else {
					// Will need to split
					sim.addStep(
						"Node Full",
						fmt.Sprintf("Leaf node is full (%d keys). Will need to split after insertion", len(node.Keys)),
						[]protocol.Highlight{{Type: "node", ID: nodeID, Color: "#ef4444", Animation: "shake"}},
						sim.cloneTreeData(treeCopy),
					)
				}
			}
		}

		// Perform actual insertion
		oldNodes := make(map[string]bool)
		for id := range treeCopy.Nodes {
			oldNodes[id] = true
		}

		treeCopy.Insert(key)

		// Check for splits
		newNodes := []string{}
		for id := range treeCopy.Nodes {
			if !oldNodes[id] {
				newNodes = append(newNodes, id)
			}
		}

		if len(newNodes) > 0 {
			highlights := []protocol.Highlight{}
			for _, id := range newNodes {
				highlights = append(highlights, protocol.Highlight{
					Type:      "node",
					ID:        id,
					Color:     "#f59e0b",
					Animation: "fadeIn",
				})
			}
			sim.addStep(
				"Split Occurred",
				fmt.Sprintf("Node split required. Created %d new node(s)", len(newNodes)),
				highlights,
				sim.cloneTreeData(treeCopy),
			)
		}

		sim.addStep(
			"Insertion Complete",
			fmt.Sprintf("Key %d successfully inserted into the B-Tree", key),
			[]protocol.Highlight{{Type: "key", ID: fmt.Sprintf("%d", key), Color: "#10b981", Animation: "pulse"}},
			sim.cloneTreeData(treeCopy),
		)
	}

	// Apply to actual tree
	sim.tree.Insert(key)
}

// PrepareSearch generates steps for a search operation
func (sim *BTreeSimulation) PrepareSearch(key int) {
	sim.operation = "search"
	sim.operand = key
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.searchPath = nil

	sim.addStep(
		fmt.Sprintf("Search for %d", key),
		fmt.Sprintf("Starting search for key %d in the B-Tree", key),
		[]protocol.Highlight{},
		sim.cloneTreeData(sim.tree),
	)

	if sim.tree.RootID == "" {
		sim.addStep(
			"Tree Empty",
			"The tree is empty. Key not found.",
			[]protocol.Highlight{},
			sim.cloneTreeData(sim.tree),
		)
		return
	}

	path := []string{}
	found := sim.searchWithSteps(sim.tree.RootID, key, &path)

	sim.searchPath = path

	if found {
		sim.addStep(
			"Key Found!",
			fmt.Sprintf("Key %d found in the B-Tree", key),
			[]protocol.Highlight{{Type: "key", ID: fmt.Sprintf("%d", key), Color: "#10b981", Animation: "pulse"}},
			sim.cloneTreeData(sim.tree),
		)
	} else {
		sim.addStep(
			"Key Not Found",
			fmt.Sprintf("Key %d does not exist in the B-Tree", key),
			[]protocol.Highlight{},
			sim.cloneTreeData(sim.tree),
		)
	}
}

func (sim *BTreeSimulation) searchWithSteps(nodeID string, key int, path *[]string) bool {
	node := sim.tree.Nodes[nodeID]
	*path = append(*path, nodeID)

	// Find position
	i := 0
	for i < len(node.Keys) && key > node.Keys[i] {
		i++
	}

	keyComparison := "greater than all keys"
	if i < len(node.Keys) {
		if key == node.Keys[i] {
			keyComparison = fmt.Sprintf("equal to key at position %d", i)
		} else {
			keyComparison = fmt.Sprintf("less than key %d at position %d", node.Keys[i], i)
		}
	}

	sim.addStep(
		fmt.Sprintf("Examine %s", nodeID),
		fmt.Sprintf("Comparing %d with keys %v. Result: %s", key, node.Keys, keyComparison),
		[]protocol.Highlight{
			{Type: "node", ID: nodeID, Color: "#3b82f6", Animation: "pulse"},
		},
		sim.cloneTreeData(sim.tree),
	)

	// Check if found
	if i < len(node.Keys) && key == node.Keys[i] {
		return true
	}

	// If leaf, not found
	if node.IsLeaf {
		return false
	}

	// Recurse
	childID := node.Children[i]
	sim.addStep(
		"Follow Child Pointer",
		fmt.Sprintf("Following child pointer %d to node %s", i, childID),
		[]protocol.Highlight{
			{Type: "edge", ID: fmt.Sprintf("%s-%s", nodeID, childID), Color: "#f59e0b", Animation: "pulse"},
		},
		sim.cloneTreeData(sim.tree),
	)

	return sim.searchWithSteps(childID, key, path)
}

// PrepareDelete generates steps for a delete operation
func (sim *BTreeSimulation) PrepareDelete(key int) {
	sim.operation = "delete"
	sim.operand = key
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.searchPath = nil

	treeCopy := sim.tree.Clone()

	sim.addStep(
		fmt.Sprintf("Delete %d", key),
		fmt.Sprintf("Starting deletion of key %d from the B-Tree", key),
		[]protocol.Highlight{},
		sim.cloneTreeData(treeCopy),
	)

	// First find the key
	nodeID, keyIndex, found := treeCopy.Search(key)
	if !found {
		sim.addStep(
			"Key Not Found",
			fmt.Sprintf("Key %d does not exist in the tree. Nothing to delete.", key),
			[]protocol.Highlight{},
			sim.cloneTreeData(treeCopy),
		)
		return
	}

	sim.addStep(
		"Key Found",
		fmt.Sprintf("Found key %d at node %s, position %d", key, nodeID, keyIndex),
		[]protocol.Highlight{{Type: "key", ID: fmt.Sprintf("%d", key), Color: "#ef4444", Animation: "pulse"}},
		sim.cloneTreeData(treeCopy),
	)

	node := treeCopy.Nodes[nodeID]
	if node.IsLeaf {
		sim.addStep(
			"Delete from Leaf",
			fmt.Sprintf("Key %d is in a leaf node. Removing directly.", key),
			[]protocol.Highlight{{Type: "node", ID: nodeID, Color: "#f59e0b", Animation: "pulse"}},
			sim.cloneTreeData(treeCopy),
		)
	} else {
		sim.addStep(
			"Delete from Internal",
			fmt.Sprintf("Key %d is in an internal node. Will replace with predecessor/successor.", key),
			[]protocol.Highlight{{Type: "node", ID: nodeID, Color: "#f59e0b", Animation: "pulse"}},
			sim.cloneTreeData(treeCopy),
		)
	}

	// Perform deletion
	treeCopy.Delete(key)

	sim.addStep(
		"Deletion Complete",
		fmt.Sprintf("Key %d successfully deleted from the B-Tree", key),
		[]protocol.Highlight{},
		sim.cloneTreeData(treeCopy),
	)

	// Apply to actual tree
	sim.tree.Delete(key)
}

// PrepareRangeSearch generates steps for a range search
func (sim *BTreeSimulation) PrepareRangeSearch(start, end int) {
	sim.operation = "range"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.searchPath = nil

	sim.addStep(
		fmt.Sprintf("Range [%d, %d]", start, end),
		fmt.Sprintf("Starting range search for keys between %d and %d", start, end),
		[]protocol.Highlight{},
		sim.cloneTreeData(sim.tree),
	)

	results := sim.tree.RangeSearch(start, end)

	if len(results) > 0 {
		highlights := []protocol.Highlight{}
		for _, k := range results {
			highlights = append(highlights, protocol.Highlight{
				Type:  "key",
				ID:    fmt.Sprintf("%d", k),
				Color: "#10b981",
			})
		}
		sim.addStep(
			"Range Search Complete",
			fmt.Sprintf("Found %d keys in range: %v", len(results), results),
			highlights,
			sim.cloneTreeData(sim.tree),
		)
	} else {
		sim.addStep(
			"No Results",
			fmt.Sprintf("No keys found in range [%d, %d]", start, end),
			[]protocol.Highlight{},
			sim.cloneTreeData(sim.tree),
		)
	}
}

// Helper methods
func (sim *BTreeSimulation) addStep(title, description string, highlights []protocol.Highlight, _ map[string]interface{}) {
	step := engine.Step{
		Index:       len(sim.steps),
		Title:       title,
		Description: description,
		Highlights:  highlights,
	}
	sim.steps = append(sim.steps, step)
}

func (sim *BTreeSimulation) cloneTreeData(tree *internal.BTree) map[string]interface{} {
	nodes := make(map[string]interface{})
	for id, node := range tree.Nodes {
		nodes[id] = map[string]interface{}{
			"id":       node.ID,
			"keys":     append([]int{}, node.Keys...),
			"children": append([]string{}, node.Children...),
			"isLeaf":   node.IsLeaf,
			"parent":   node.Parent,
		}
	}
	return map[string]interface{}{
		"nodes":  nodes,
		"rootId": tree.RootID,
		"order":  tree.Order,
		"path":   sim.searchPath,
	}
}

func (sim *BTreeSimulation) findInsertionPath(tree *internal.BTree, key int) []string {
	path := []string{}
	if tree.RootID == "" {
		return path
	}

	nodeID := tree.RootID
	for {
		path = append(path, nodeID)
		node := tree.Nodes[nodeID]
		if node.IsLeaf {
			break
		}
		// Find child to follow
		i := 0
		for i < len(node.Keys) && key > node.Keys[i] {
			i++
		}
		nodeID = node.Children[i]
	}
	return path
}
