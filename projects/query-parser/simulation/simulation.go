package simulation

import (
	"fmt"

	"github.com/ersantana/db-internals/packages/protocol"
	"github.com/ersantana/db-internals/packages/simulation/engine"
	"github.com/ersantana/db-internals/projects/query-parser/internal"
)

// ParserSimulation implements the simulation.Simulation interface
type ParserSimulation struct {
	query             string
	tokens            []internal.Token
	astNodes          map[string]*internal.ASTNode
	astRoot           string
	steps             []engine.Step
	currentStep       int
	currentTokenIndex int
	parsePhase        string
}

// NewParserSimulation creates a new parser simulation
func NewParserSimulation() *ParserSimulation {
	return &ParserSimulation{
		tokens:      []internal.Token{},
		astNodes:    make(map[string]*internal.ASTNode),
		steps:       make([]engine.Step, 0),
		currentStep: -1,
		parsePhase:  "idle",
	}
}

// Name returns the simulation name
func (sim *ParserSimulation) Name() string {
	return "Query Parser"
}

// Description returns the simulation description
func (sim *ParserSimulation) Description() string {
	return "Interactive visualization of SQL query parsing"
}

// Initialize sets up the simulation with given config
func (sim *ParserSimulation) Initialize(config map[string]interface{}) error {
	sim.tokens = []internal.Token{}
	sim.astNodes = make(map[string]*internal.ASTNode)
	sim.astRoot = ""
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.currentTokenIndex = -1
	sim.parsePhase = "idle"

	// If query is provided, prepare the simulation
	if query, ok := config["query"].(string); ok && query != "" {
		sim.query = query
		sim.prepareParseSimulation(query)
	}

	return nil
}

// Reset returns the simulation to initial state
func (sim *ParserSimulation) Reset() error {
	sim.tokens = []internal.Token{}
	sim.astNodes = make(map[string]*internal.ASTNode)
	sim.astRoot = ""
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	sim.currentTokenIndex = -1
	sim.parsePhase = "idle"
	return nil
}

// GenerateSteps returns all steps for current simulation
func (sim *ParserSimulation) GenerateSteps() []engine.Step {
	return sim.steps
}

// CurrentStep returns current step index
func (sim *ParserSimulation) CurrentStep() int {
	return sim.currentStep
}

// ExecuteStep executes a specific step
func (sim *ParserSimulation) ExecuteStep(index int) engine.StepResult {
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
func (sim *ParserSimulation) CanStepForward() bool {
	return sim.currentStep < len(sim.steps)-1
}

// CanStepBackward returns true if can go back
func (sim *ParserSimulation) CanStepBackward() bool {
	return sim.currentStep > 0
}

// GetState returns current state
func (sim *ParserSimulation) GetState() interface{} {
	return map[string]interface{}{
		"query":             sim.query,
		"tokens":            sim.tokens,
		"astNodes":          sim.astNodes,
		"astRoot":           sim.astRoot,
		"currentStep":       sim.currentStep,
		"currentTokenIndex": sim.currentTokenIndex,
		"parsePhase":        sim.parsePhase,
	}
}

// GetVisualizationData returns data for rendering
func (sim *ParserSimulation) GetVisualizationData() map[string]interface{} {
	// Convert tokens to interface slice
	tokenData := make([]map[string]interface{}, len(sim.tokens))
	for i, t := range sim.tokens {
		tokenData[i] = map[string]interface{}{
			"type":     t.Type,
			"value":    t.Value,
			"position": t.Position,
		}
	}

	// Convert AST nodes
	nodeData := make(map[string]interface{})
	for id, node := range sim.astNodes {
		nodeData[id] = map[string]interface{}{
			"id":       node.ID,
			"type":     node.Type,
			"value":    node.Value,
			"children": node.Children,
			"parent":   node.Parent,
			"meta":     node.Meta,
		}
	}

	return map[string]interface{}{
		"query":             sim.query,
		"tokens":            tokenData,
		"astNodes":          nodeData,
		"astRoot":           sim.astRoot,
		"currentTokenIndex": sim.currentTokenIndex,
		"parsePhase":        sim.parsePhase,
	}
}

// prepareParseSimulation generates all steps for parsing a query
func (sim *ParserSimulation) prepareParseSimulation(query string) {
	sim.query = query
	sim.steps = make([]engine.Step, 0)
	sim.parsePhase = "tokenizing"

	// Step: Start
	sim.addStep(
		"Start Parsing",
		fmt.Sprintf("Parsing SQL query: %s", query),
		[]protocol.Highlight{},
	)

	// Tokenization phase
	sim.addStep(
		"Tokenization",
		"Breaking input into tokens (lexical analysis)",
		[]protocol.Highlight{},
	)

	lexer := internal.NewLexer(query)

	for lexer.HasMore() {
		token := lexer.TokenizeStep()
		if token.Type == internal.TokenEOF {
			break
		}

		sim.tokens = append(sim.tokens, *token)
		tokenIndex := len(sim.tokens) - 1

		sim.addStep(
			fmt.Sprintf("Token: %s", token.Type),
			fmt.Sprintf("Found %s token: '%s' at position %d-%d",
				token.Type, token.Value, token.Position.Start, token.Position.End),
			[]protocol.Highlight{
				{Type: "token", ID: fmt.Sprintf("token-%d", tokenIndex), Color: "#3b82f6", Animation: "pulse"},
			},
		)
	}

	// Add EOF token
	eofToken := internal.Token{
		Type:     internal.TokenEOF,
		Value:    "",
		Position: internal.Position{Start: len(query), End: len(query)},
	}
	sim.tokens = append(sim.tokens, eofToken)

	sim.addStep(
		"Tokenization Complete",
		fmt.Sprintf("Created %d tokens from input", len(sim.tokens)-1),
		[]protocol.Highlight{},
	)

	// Parsing phase
	sim.parsePhase = "parsing"
	sim.addStep(
		"Parsing",
		"Building Abstract Syntax Tree (AST) from tokens",
		[]protocol.Highlight{},
	)

	parser := internal.NewParser(sim.tokens)
	root, err := parser.Parse()

	if err != nil {
		sim.addStep(
			"Parse Error",
			err.Error(),
			[]protocol.Highlight{},
		)
		sim.parsePhase = "error"
		return
	}

	sim.astNodes = parser.GetNodes()
	sim.astRoot = parser.GetRootID()

	// Generate steps for each AST node
	sim.generateASTSteps(root)

	sim.parsePhase = "complete"
	sim.addStep(
		"Parsing Complete",
		fmt.Sprintf("Successfully created AST with %d nodes", len(sim.astNodes)),
		[]protocol.Highlight{
			{Type: "node", ID: sim.astRoot, Color: "#10b981", Animation: "pulse"},
		},
	)
}

func (sim *ParserSimulation) generateASTSteps(node *internal.ASTNode) {
	if node == nil {
		return
	}

	nodeType := string(node.Type)
	description := fmt.Sprintf("Created %s node", nodeType)
	if node.Value != "" {
		description = fmt.Sprintf("Created %s node with value '%s'", nodeType, node.Value)
	}

	sim.addStep(
		fmt.Sprintf("AST: %s", nodeType),
		description,
		[]protocol.Highlight{
			{Type: "node", ID: node.ID, Color: "#8b5cf6", Animation: "fadeIn"},
		},
	)

	// Recurse for children
	for _, childID := range node.Children {
		if childNode, ok := sim.astNodes[childID]; ok {
			sim.generateASTSteps(childNode)
		}
	}
}

// ParseQuery starts parsing a new query
func (sim *ParserSimulation) ParseQuery(query string) {
	sim.Reset()
	sim.prepareParseSimulation(query)
}

// Helper method
func (sim *ParserSimulation) addStep(title, description string, highlights []protocol.Highlight) {
	step := engine.Step{
		Index:       len(sim.steps),
		Title:       title,
		Description: description,
		Highlights:  highlights,
	}
	sim.steps = append(sim.steps, step)
}
