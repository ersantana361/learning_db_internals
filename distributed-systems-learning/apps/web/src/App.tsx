import { useEffect } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { useSimulationStore } from './stores/simulationStore';
import { ControlPanel } from './components/common/ControlPanel';
import { NodeVisualizer } from './components/common/NodeVisualizer';
import { Timeline } from './components/common/Timeline';
import { ProjectSelector } from './components/common/ProjectSelector';
import './App.css';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

function App() {
  const { currentProject, setCurrentProject, simulation } = useSimulationStore();

  const {
    isConnected,
    error,
    startSimulation,
    pauseSimulation,
    resumeSimulation,
    stopSimulation,
    stepForward,
    setSpeed,
    injectCrash,
    recoverNode,
    getState,
  } = useWebSocket({
    url: WS_URL,
    onOpen: () => {
      console.log('Connected to simulation server');
      getState();
    },
    onClose: () => {
      console.log('Disconnected from simulation server');
    },
  });

  useEffect(() => {
    // Poll for state updates periodically
    const interval = setInterval(() => {
      if (isConnected && simulation.running) {
        getState();
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [isConnected, simulation.running, getState]);

  const handleProjectSelect = (projectId: string) => {
    setCurrentProject(projectId);
  };

  const handleStart = () => {
    if (currentProject) {
      startSimulation(currentProject, undefined, {
        nodeCount: 5,
        speed: 1.0,
        stepMode: true,
      });
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>Distributed Systems Learning Platform</h1>
        <div className="connection-status">
          <span className={`status-dot ${isConnected ? 'connected' : 'disconnected'}`} />
          {isConnected ? 'Connected' : 'Connecting...'}
        </div>
      </header>

      {error && (
        <div className="error-banner">
          {error}
        </div>
      )}

      <main className="app-main">
        {!currentProject ? (
          <ProjectSelector
            onSelect={handleProjectSelect}
            currentProject={currentProject}
          />
        ) : (
          <div className="simulation-view">
            <div className="simulation-header">
              <button
                className="back-button"
                onClick={() => {
                  stopSimulation();
                  setCurrentProject(null);
                }}
              >
                Back to Projects
              </button>
              <h2>{currentProject.replace(/-/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}</h2>
            </div>

            <ControlPanel
              onStart={handleStart}
              onPause={pauseSimulation}
              onResume={resumeSimulation}
              onStop={stopSimulation}
              onStep={stepForward}
              onSpeedChange={setSpeed}
            />

            <div className="simulation-content">
              <div className="visualization-panel">
                <NodeVisualizer
                  onNodeClick={(nodeId) => console.log('Clicked node:', nodeId)}
                  onCrashNode={injectCrash}
                  onRecoverNode={recoverNode}
                />
              </div>

              <div className="timeline-panel">
                <Timeline />
              </div>
            </div>

            <div className="explanation-panel">
              <h3>About This Project</h3>
              <p>
                This visualization helps you understand distributed systems concepts
                through interactive simulations. Use the controls above to start,
                pause, and step through the simulation. Click on nodes to inspect
                their state or inject failures.
              </p>
            </div>
          </div>
        )}
      </main>

      <footer className="app-footer">
        <p>Built for learning distributed systems concepts from Martin Kleppmann's course</p>
      </footer>
    </div>
  );
}

export default App;
