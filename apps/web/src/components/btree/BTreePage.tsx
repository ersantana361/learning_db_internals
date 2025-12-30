import { useState } from 'react';
import { ArrowLeft, Plus, Search, Trash2, Play } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useWebSocket } from '../../hooks/useWebSocket';
import { SimulationControls } from '../common/SimulationControls';
import { StepTimeline } from '../common/StepTimeline';
import { BTreeVisualization } from './BTreeVisualization';
import './BTreePage.css';

const SCENARIOS = [
  { id: 'basic-insert', name: 'Basic Insertion', description: 'Insert keys without splits' },
  { id: 'split-demo', name: 'Node Splitting', description: 'Demonstrates node splits' },
  { id: 'search-demo', name: 'Search Operations', description: 'Search traversal' },
  { id: 'delete-demo', name: 'Delete Operations', description: 'Deletion and rebalancing' },
  { id: 'large-tree', name: 'Multi-Level Tree', description: 'Build a larger tree' },
];

export function BTreePage() {
  const [inputValue, setInputValue] = useState('');
  const [operation, setOperation] = useState<'insert' | 'search' | 'delete'>('insert');
  const [selectedScenario, setSelectedScenario] = useState('');

  const {
    isConnected,
    isConnecting,
    startSimulation,
    stepForward,
    stepBackward,
    play,
    pause,
    reset,
    setSpeed,
    executeOperation,
    selectScenario,
  } = useWebSocket();

  const handleStartSimulation = () => {
    startSimulation({
      project: 'btree',
      parameters: {
        order: 4,
      },
    });
  };

  const handleExecuteOperation = () => {
    const value = parseInt(inputValue, 10);
    if (isNaN(value)) return;

    executeOperation(operation, { key: value });
    setInputValue('');
  };

  const handleScenarioSelect = (scenarioId: string) => {
    setSelectedScenario(scenarioId);
    selectScenario(scenarioId);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleExecuteOperation();
    }
  };

  return (
    <div className="btree-page">
      <header className="page-header">
        <Link to="/" className="back-link">
          <ArrowLeft size={20} />
          Back to Topics
        </Link>
        <div className="header-content">
          <h1>B-Tree Visualization</h1>
          <p>Interactive exploration of B-Tree operations</p>
        </div>
        <div className="connection-badge">
          <span className={`status-dot ${isConnected ? 'connected' : 'disconnected'}`} />
          {isConnecting ? 'Connecting...' : isConnected ? 'Connected' : 'Disconnected'}
        </div>
      </header>

      <div className="page-content">
        <aside className="control-panel">
          <section className="panel-section">
            <h3>Quick Start</h3>
            <button
              className="start-btn"
              onClick={handleStartSimulation}
              disabled={!isConnected}
            >
              <Play size={16} />
              Start New Simulation
            </button>
          </section>

          <section className="panel-section">
            <h3>Scenarios</h3>
            <div className="scenario-list">
              {SCENARIOS.map((scenario) => (
                <button
                  key={scenario.id}
                  className={`scenario-btn ${selectedScenario === scenario.id ? 'selected' : ''}`}
                  onClick={() => handleScenarioSelect(scenario.id)}
                  disabled={!isConnected}
                >
                  <span className="scenario-name">{scenario.name}</span>
                  <span className="scenario-desc">{scenario.description}</span>
                </button>
              ))}
            </div>
          </section>

          <section className="panel-section">
            <h3>Manual Operation</h3>
            <div className="operation-selector">
              <button
                className={`op-btn ${operation === 'insert' ? 'active' : ''}`}
                onClick={() => setOperation('insert')}
              >
                <Plus size={14} />
                Insert
              </button>
              <button
                className={`op-btn ${operation === 'search' ? 'active' : ''}`}
                onClick={() => setOperation('search')}
              >
                <Search size={14} />
                Search
              </button>
              <button
                className={`op-btn ${operation === 'delete' ? 'active' : ''}`}
                onClick={() => setOperation('delete')}
              >
                <Trash2 size={14} />
                Delete
              </button>
            </div>
            <div className="input-group">
              <input
                type="number"
                placeholder="Enter key value"
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyPress={handleKeyPress}
                disabled={!isConnected}
              />
              <button
                className="execute-btn"
                onClick={handleExecuteOperation}
                disabled={!isConnected || !inputValue}
              >
                Execute
              </button>
            </div>
          </section>

          <section className="panel-section">
            <h3>Playback Controls</h3>
            <SimulationControls
              onStepForward={stepForward}
              onStepBackward={stepBackward}
              onPlay={play}
              onPause={pause}
              onReset={reset}
              onSpeedChange={setSpeed}
              disabled={!isConnected}
            />
          </section>

          <section className="panel-section">
            <h3>Step Timeline</h3>
            <StepTimeline />
          </section>
        </aside>

        <main className="visualization-area">
          <BTreeVisualization />
        </main>
      </div>
    </div>
  );
}

export default BTreePage;
