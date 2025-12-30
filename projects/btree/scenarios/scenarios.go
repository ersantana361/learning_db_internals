package scenarios

// Scenario represents a predefined B-Tree scenario
type Scenario struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Operations  []Operation            `json:"operations"`
}

// Operation represents an operation to perform
type Operation struct {
	Type   string                 `json:"type"` // insert, search, delete, range
	Params map[string]interface{} `json:"params"`
}

// GetScenarios returns all available scenarios
func GetScenarios() []Scenario {
	return []Scenario{
		BasicInsert(),
		SplitDemo(),
		SearchDemo(),
		DeleteDemo(),
		RangeQueryDemo(),
		LargeTreeDemo(),
	}
}

// GetScenario returns a scenario by ID
func GetScenario(id string) *Scenario {
	for _, s := range GetScenarios() {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

// BasicInsert demonstrates simple insertions without splits
func BasicInsert() Scenario {
	return Scenario{
		ID:          "basic-insert",
		Name:        "Basic Insertion",
		Description: "Insert keys into an empty B-Tree without triggering splits",
		Config: map[string]interface{}{
			"order": 4,
		},
		Operations: []Operation{
			{Type: "insert", Params: map[string]interface{}{"key": 10}},
			{Type: "insert", Params: map[string]interface{}{"key": 20}},
			{Type: "insert", Params: map[string]interface{}{"key": 5}},
		},
	}
}

// SplitDemo demonstrates node splitting during insertion
func SplitDemo() Scenario {
	return Scenario{
		ID:          "split-demo",
		Name:        "Node Splitting",
		Description: "Demonstrates how nodes split when they become full",
		Config: map[string]interface{}{
			"order":       4,
			"initialKeys": []int{10, 20, 30},
		},
		Operations: []Operation{
			{Type: "insert", Params: map[string]interface{}{"key": 40}},
			{Type: "insert", Params: map[string]interface{}{"key": 50}},
			{Type: "insert", Params: map[string]interface{}{"key": 25}},
			{Type: "insert", Params: map[string]interface{}{"key": 35}},
		},
	}
}

// SearchDemo demonstrates search operations
func SearchDemo() Scenario {
	return Scenario{
		ID:          "search-demo",
		Name:        "Search Operations",
		Description: "Demonstrates how search traverses the B-Tree",
		Config: map[string]interface{}{
			"order":       4,
			"initialKeys": []int{10, 20, 30, 40, 50, 60, 70, 80},
		},
		Operations: []Operation{
			{Type: "search", Params: map[string]interface{}{"key": 50}},
			{Type: "search", Params: map[string]interface{}{"key": 25}},
			{Type: "search", Params: map[string]interface{}{"key": 70}},
		},
	}
}

// DeleteDemo demonstrates delete operations
func DeleteDemo() Scenario {
	return Scenario{
		ID:          "delete-demo",
		Name:        "Delete Operations",
		Description: "Demonstrates key deletion and node rebalancing",
		Config: map[string]interface{}{
			"order":       4,
			"initialKeys": []int{10, 20, 30, 40, 50, 60, 70},
		},
		Operations: []Operation{
			{Type: "delete", Params: map[string]interface{}{"key": 30}},
			{Type: "delete", Params: map[string]interface{}{"key": 50}},
			{Type: "delete", Params: map[string]interface{}{"key": 10}},
		},
	}
}

// RangeQueryDemo demonstrates range queries
func RangeQueryDemo() Scenario {
	return Scenario{
		ID:          "range-query",
		Name:        "Range Queries",
		Description: "Demonstrates range search operations",
		Config: map[string]interface{}{
			"order":       4,
			"initialKeys": []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50},
		},
		Operations: []Operation{
			{Type: "range", Params: map[string]interface{}{"start": 15, "end": 35}},
			{Type: "range", Params: map[string]interface{}{"start": 1, "end": 10}},
			{Type: "range", Params: map[string]interface{}{"start": 100, "end": 200}},
		},
	}
}

// LargeTreeDemo demonstrates a larger tree with multiple levels
func LargeTreeDemo() Scenario {
	return Scenario{
		ID:          "large-tree",
		Name:        "Multi-Level Tree",
		Description: "Build a larger B-Tree with multiple levels",
		Config: map[string]interface{}{
			"order": 4,
		},
		Operations: []Operation{
			{Type: "insert", Params: map[string]interface{}{"key": 50}},
			{Type: "insert", Params: map[string]interface{}{"key": 25}},
			{Type: "insert", Params: map[string]interface{}{"key": 75}},
			{Type: "insert", Params: map[string]interface{}{"key": 10}},
			{Type: "insert", Params: map[string]interface{}{"key": 30}},
			{Type: "insert", Params: map[string]interface{}{"key": 60}},
			{Type: "insert", Params: map[string]interface{}{"key": 90}},
			{Type: "insert", Params: map[string]interface{}{"key": 5}},
			{Type: "insert", Params: map[string]interface{}{"key": 15}},
			{Type: "insert", Params: map[string]interface{}{"key": 27}},
			{Type: "insert", Params: map[string]interface{}{"key": 35}},
			{Type: "insert", Params: map[string]interface{}{"key": 55}},
			{Type: "insert", Params: map[string]interface{}{"key": 65}},
			{Type: "insert", Params: map[string]interface{}{"key": 80}},
			{Type: "insert", Params: map[string]interface{}{"key": 95}},
		},
	}
}
