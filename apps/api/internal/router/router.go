package router

import (
	"encoding/json"
	"net/http"

	"github.com/ersantana/db-internals/apps/api/internal/handlers"
)

// New creates a new HTTP router with the simulation manager
func New(hub *handlers.Hub, simManager *handlers.SimulationManager) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// WebSocket endpoint with simulation manager routing
	mux.HandleFunc("/ws", handlers.HandleWebSocket(
		hub,
		simManager.HandleMessage,
		simManager.RemoveSession,
	))

	// API endpoints
	mux.HandleFunc("/api/topics", handleTopics)
	mux.HandleFunc("/api/projects", handleProjects(simManager))

	return mux
}

// Topic represents a learning topic
type Topic struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	topics := []Topic{
		{ID: "btree", Name: "B-Tree", Description: "Self-balancing tree data structure for sorted data", Category: "storage"},
		{ID: "lsm-tree", Name: "LSM Tree", Description: "Log-structured merge tree for write-optimized storage", Category: "storage"},
		{ID: "buffer-pool", Name: "Buffer Pool", Description: "In-memory page cache management", Category: "storage"},
		{ID: "page-layout", Name: "Page Layout", Description: "On-disk page structure and organization", Category: "storage"},
		{ID: "query-parser", Name: "Query Parser", Description: "SQL parsing and AST generation", Category: "query"},
		{ID: "query-optimizer", Name: "Query Optimizer", Description: "Query plan optimization and cost estimation", Category: "query"},
		{ID: "execution-engine", Name: "Execution Engine", Description: "Query plan execution and operators", Category: "query"},
		{ID: "mvcc", Name: "MVCC", Description: "Multi-version concurrency control", Category: "transaction"},
		{ID: "wal", Name: "Write-Ahead Log", Description: "Durability and crash recovery", Category: "transaction"},
		{ID: "locking", Name: "Locking", Description: "Lock management and deadlock detection", Category: "transaction"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

func handleProjects(simManager *handlers.SimulationManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects := simManager.GetRegisteredProjects()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	}
}
