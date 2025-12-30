package simulation

import (
	"fmt"

	"github.com/ersantana/db-internals/packages/protocol"
	"github.com/ersantana/db-internals/packages/simulation/engine"
	"github.com/ersantana/db-internals/projects/mvcc/internal"
)

// MVCCSimulation implements the simulation.Simulation interface
type MVCCSimulation struct {
	store       *internal.MVCCStore
	steps       []engine.Step
	currentStep int
	operation   string
}

// NewMVCCSimulation creates a new MVCC simulation
func NewMVCCSimulation() *MVCCSimulation {
	return &MVCCSimulation{
		store:       internal.NewMVCCStore(),
		steps:       make([]engine.Step, 0),
		currentStep: -1,
	}
}

// Name returns the simulation name
func (sim *MVCCSimulation) Name() string {
	return "MVCC"
}

// Description returns the simulation description
func (sim *MVCCSimulation) Description() string {
	return "Interactive visualization of Multi-Version Concurrency Control"
}

// Initialize sets up the simulation with given config
func (sim *MVCCSimulation) Initialize(config map[string]interface{}) error {
	sim.store = internal.NewMVCCStore()

	// Pre-populate with initial data if specified
	if populate, ok := config["initialData"].(bool); ok && populate {
		sim.store.InsertInitialData()
	}

	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	return nil
}

// Reset returns the simulation to initial state
func (sim *MVCCSimulation) Reset() error {
	sim.store = internal.NewMVCCStore()
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1
	return nil
}

// GenerateSteps returns all steps for current simulation
func (sim *MVCCSimulation) GenerateSteps() []engine.Step {
	return sim.steps
}

// CurrentStep returns current step index
func (sim *MVCCSimulation) CurrentStep() int {
	return sim.currentStep
}

// ExecuteStep executes a specific step
func (sim *MVCCSimulation) ExecuteStep(index int) engine.StepResult {
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
func (sim *MVCCSimulation) CanStepForward() bool {
	return sim.currentStep < len(sim.steps)-1
}

// CanStepBackward returns true if can go back
func (sim *MVCCSimulation) CanStepBackward() bool {
	return sim.currentStep > 0
}

// GetState returns current store state
func (sim *MVCCSimulation) GetState() interface{} {
	return map[string]interface{}{
		"transactions":    sim.store.Transactions,
		"versions":        sim.store.Versions,
		"rows":            sim.store.Rows,
		"globalTimestamp": sim.store.GlobalTimestamp,
		"currentStep":     sim.currentStep,
		"totalSteps":      len(sim.steps),
	}
}

// GetVisualizationData returns data for rendering
func (sim *MVCCSimulation) GetVisualizationData() map[string]interface{} {
	return map[string]interface{}{
		"transactions":      sim.store.Transactions,
		"versions":          sim.store.Versions,
		"rows":              sim.store.Rows,
		"globalTimestamp":   sim.store.GlobalTimestamp,
		"activeTransaction": sim.store.ActiveTransaction,
	}
}

// PrepareBeginTransaction generates steps for beginning a transaction
func (sim *MVCCSimulation) PrepareBeginTransaction() {
	sim.operation = "begin"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1

	storeCopy := sim.store.Clone()

	sim.addStep(
		"Begin Transaction",
		fmt.Sprintf("Starting new transaction at timestamp %d", storeCopy.GlobalTimestamp),
		[]protocol.Highlight{},
	)

	tx := storeCopy.BeginTransaction()

	sim.addStep(
		"Transaction Started",
		fmt.Sprintf("Created transaction %s with start time %d", tx.ID, tx.StartTime),
		[]protocol.Highlight{{Type: "row", ID: tx.ID, Color: "#10b981", Animation: "pulse"}},
	)

	// Apply to actual store
	sim.store.BeginTransaction()
}

// PrepareRead generates steps for a read operation
func (sim *MVCCSimulation) PrepareRead(txID, rowID string) {
	sim.operation = "read"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1

	tx := sim.store.Transactions[txID]
	if tx == nil {
		sim.addStep(
			"Error",
			fmt.Sprintf("Transaction %s not found", txID),
			[]protocol.Highlight{},
		)
		return
	}

	row := sim.store.Rows[rowID]
	if row == nil {
		sim.addStep(
			"Error",
			fmt.Sprintf("Row %s not found", rowID),
			[]protocol.Highlight{},
		)
		return
	}

	sim.addStep(
		fmt.Sprintf("Read %s", rowID),
		fmt.Sprintf("Transaction %s reading row %s at snapshot time %d", txID, rowID, tx.StartTime),
		[]protocol.Highlight{
			{Type: "row", ID: txID, Color: "#3b82f6", Animation: "pulse"},
		},
	)

	// Walk through version chain
	for _, verID := range row.VersionChain {
		ver := sim.store.Versions[verID]
		creatorTx := sim.store.Transactions[ver.CreatedBy]

		visible := false
		reason := ""

		if ver.CreatedBy == txID {
			visible = true
			reason = "Created by current transaction"
		} else if creatorTx.Status != internal.TxCommitted {
			reason = fmt.Sprintf("Creator transaction %s not committed", ver.CreatedBy)
		} else if creatorTx.CommitTime == nil || *creatorTx.CommitTime > tx.StartTime {
			reason = fmt.Sprintf("Committed at %d, after our snapshot %d", *creatorTx.CommitTime, tx.StartTime)
		} else {
			visible = true
			reason = fmt.Sprintf("Committed at %d, before our snapshot %d", *creatorTx.CommitTime, tx.StartTime)
		}

		color := "#ef4444" // Not visible
		if visible {
			color = "#10b981"
		}

		sim.addStep(
			fmt.Sprintf("Check %s", verID),
			reason,
			[]protocol.Highlight{
				{Type: "cell", ID: verID, Color: color, Animation: "pulse"},
			},
		)

		if visible {
			sim.addStep(
				"Version Found",
				fmt.Sprintf("Reading version %s with data: %v", verID, ver.Data),
				[]protocol.Highlight{
					{Type: "cell", ID: verID, Color: "#10b981", Animation: "pulse"},
				},
			)
			break
		}
	}

	// Apply the read
	sim.store.Read(txID, rowID)
}

// PrepareWrite generates steps for a write operation
func (sim *MVCCSimulation) PrepareWrite(txID, rowID string, data map[string]interface{}) {
	sim.operation = "write"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1

	tx := sim.store.Transactions[txID]
	if tx == nil {
		sim.addStep(
			"Error",
			fmt.Sprintf("Transaction %s not found", txID),
			[]protocol.Highlight{},
		)
		return
	}

	sim.addStep(
		fmt.Sprintf("Write to %s", rowID),
		fmt.Sprintf("Transaction %s writing to row %s", txID, rowID),
		[]protocol.Highlight{
			{Type: "row", ID: txID, Color: "#f59e0b", Animation: "pulse"},
		},
	)

	row, exists := sim.store.Rows[rowID]
	if !exists {
		sim.addStep(
			"Create Row",
			fmt.Sprintf("Row %s doesn't exist, creating new row", rowID),
			[]protocol.Highlight{},
		)
	} else {
		sim.addStep(
			"Create Version",
			fmt.Sprintf("Adding new version to row %s's version chain", rowID),
			[]protocol.Highlight{
				{Type: "row", ID: row.CurrentVersion, Color: "#94a3b8"},
			},
		)
	}

	// Apply the write
	ver, _ := sim.store.Write(txID, rowID, data)

	sim.addStep(
		"Version Created",
		fmt.Sprintf("Created version %s with data: %v", ver.ID, data),
		[]protocol.Highlight{
			{Type: "cell", ID: ver.ID, Color: "#10b981", Animation: "fadeIn"},
		},
	)
}

// PrepareCommit generates steps for committing a transaction
func (sim *MVCCSimulation) PrepareCommit(txID string) {
	sim.operation = "commit"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1

	tx := sim.store.Transactions[txID]
	if tx == nil {
		sim.addStep(
			"Error",
			fmt.Sprintf("Transaction %s not found", txID),
			[]protocol.Highlight{},
		)
		return
	}

	sim.addStep(
		fmt.Sprintf("Commit %s", txID),
		fmt.Sprintf("Committing transaction %s with %d writes", txID, len(tx.WriteSet)),
		[]protocol.Highlight{
			{Type: "row", ID: txID, Color: "#f59e0b", Animation: "pulse"},
		},
	)

	// Show all versions that will become visible
	for _, rowID := range tx.WriteSet {
		row := sim.store.Rows[rowID]
		if row != nil && len(row.VersionChain) > 0 {
			sim.addStep(
				fmt.Sprintf("Finalize %s", rowID),
				fmt.Sprintf("Version for row %s will now be visible to new transactions", rowID),
				[]protocol.Highlight{
					{Type: "cell", ID: row.VersionChain[0], Color: "#10b981", Animation: "pulse"},
				},
			)
		}
	}

	// Apply the commit
	sim.store.Commit(txID)
	tx = sim.store.Transactions[txID]

	sim.addStep(
		"Commit Complete",
		fmt.Sprintf("Transaction %s committed at timestamp %d", txID, *tx.CommitTime),
		[]protocol.Highlight{
			{Type: "row", ID: txID, Color: "#10b981", Animation: "pulse"},
		},
	)
}

// PrepareAbort generates steps for aborting a transaction
func (sim *MVCCSimulation) PrepareAbort(txID string) {
	sim.operation = "abort"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1

	tx := sim.store.Transactions[txID]
	if tx == nil {
		sim.addStep(
			"Error",
			fmt.Sprintf("Transaction %s not found", txID),
			[]protocol.Highlight{},
		)
		return
	}

	sim.addStep(
		fmt.Sprintf("Abort %s", txID),
		fmt.Sprintf("Aborting transaction %s", txID),
		[]protocol.Highlight{
			{Type: "row", ID: txID, Color: "#ef4444", Animation: "shake"},
		},
	)

	// Show all versions that will be removed
	for _, rowID := range tx.WriteSet {
		row := sim.store.Rows[rowID]
		if row != nil {
			for _, verID := range row.VersionChain {
				ver := sim.store.Versions[verID]
				if ver.CreatedBy == txID {
					sim.addStep(
						fmt.Sprintf("Remove %s", verID),
						fmt.Sprintf("Removing uncommitted version %s", verID),
						[]protocol.Highlight{
							{Type: "cell", ID: verID, Color: "#ef4444", Animation: "pulse"},
						},
					)
				}
			}
		}
	}

	// Apply the abort
	sim.store.Abort(txID)

	sim.addStep(
		"Abort Complete",
		fmt.Sprintf("Transaction %s aborted, all changes rolled back", txID),
		[]protocol.Highlight{
			{Type: "row", ID: txID, Color: "#ef4444"},
		},
	)
}

// PrepareGarbageCollect generates steps for garbage collection
func (sim *MVCCSimulation) PrepareGarbageCollect() {
	sim.operation = "gc"
	sim.steps = make([]engine.Step, 0)
	sim.currentStep = -1

	sim.addStep(
		"Garbage Collection",
		"Starting garbage collection to remove old versions",
		[]protocol.Highlight{},
	)

	// Find oldest active transaction
	oldestActive := sim.store.GlobalTimestamp
	for _, tx := range sim.store.Transactions {
		if tx.Status == internal.TxActive && tx.StartTime < oldestActive {
			oldestActive = tx.StartTime
		}
	}

	sim.addStep(
		"Find Threshold",
		fmt.Sprintf("Oldest active transaction started at %d, can remove versions older than this", oldestActive),
		[]protocol.Highlight{},
	)

	removed := sim.store.GarbageCollect()

	if len(removed) > 0 {
		for _, verID := range removed {
			sim.addStep(
				fmt.Sprintf("Remove %s", verID),
				fmt.Sprintf("Removed old version %s", verID),
				[]protocol.Highlight{
					{Type: "cell", ID: verID, Color: "#ef4444", Animation: "pulse"},
				},
			)
		}
	}

	sim.addStep(
		"GC Complete",
		fmt.Sprintf("Garbage collection complete, removed %d versions", len(removed)),
		[]protocol.Highlight{},
	)
}

// Helper method
func (sim *MVCCSimulation) addStep(title, description string, highlights []protocol.Highlight) {
	step := engine.Step{
		Index:       len(sim.steps),
		Title:       title,
		Description: description,
		Highlights:  highlights,
	}
	sim.steps = append(sim.steps, step)
}
