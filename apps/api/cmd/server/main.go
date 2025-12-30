package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ersantana/db-internals/apps/api/internal/handlers"
	"github.com/ersantana/db-internals/apps/api/internal/router"
	"github.com/ersantana/db-internals/packages/simulation/engine"
	btreesim "github.com/ersantana/db-internals/projects/btree/simulation"
	mvccsim "github.com/ersantana/db-internals/projects/mvcc/simulation"
	parsersim "github.com/ersantana/db-internals/projects/query-parser/simulation"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create hub for WebSocket connections
	hub := handlers.NewHub()
	go hub.Run()

	// Create simulation manager
	simManager := handlers.NewSimulationManager(hub)

	// Register simulation projects
	simManager.RegisterProject("btree", func() engine.Simulation {
		return btreesim.NewBTreeSimulation()
	})
	simManager.RegisterProject("mvcc", func() engine.Simulation {
		return mvccsim.NewMVCCSimulation()
	})
	simManager.RegisterProject("query-parser", func() engine.Simulation {
		return parsersim.NewParserSimulation()
	})

	// Create router with simulation manager
	r := router.New(hub, simManager)

	log.Printf("Starting server on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost:%s/ws", port)
	log.Printf("API endpoint: http://localhost:%s/api", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
