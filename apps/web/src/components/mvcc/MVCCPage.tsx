import { useState } from 'react';
import { ArrowLeft, Play, PlayCircle, XCircle, CheckCircle, Eye, Edit } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useWebSocket } from '../../hooks/useWebSocket';
import { useMVCCStore } from '../../stores/mvccStore';
import { SimulationControls } from '../common/SimulationControls';
import { StepTimeline } from '../common/StepTimeline';
import { TransactionTimeline } from './TransactionTimeline';
import { VersionChainView } from './VersionChainView';
import './MVCCPage.css';

const SCENARIOS = [
  { id: 'concurrent-reads', name: 'Concurrent Reads', description: 'Multiple transactions reading same data' },
  { id: 'write-conflict', name: 'Write Conflict', description: 'Conflicting writes demonstration' },
  { id: 'snapshot-isolation', name: 'Snapshot Isolation', description: 'How snapshots work' },
  { id: 'garbage-collection', name: 'Garbage Collection', description: 'Cleanup of old versions' },
];

type Operation = 'begin' | 'read' | 'write' | 'commit' | 'abort';

export function MVCCPage() {
  const [operation, setOperation] = useState<Operation>('begin');
  const [selectedTx, setSelectedTx] = useState('');
  const [rowId, setRowId] = useState('');
  const [writeValue, setWriteValue] = useState('');
  const [selectedScenario, setSelectedScenario] = useState('');

  const { transactions } = useMVCCStore();

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

  const activeTxs = Object.values(transactions).filter(tx => tx.status === 'active');

  const handleStartSimulation = () => {
    startSimulation({
      project: 'mvcc',
      parameters: {
        initialData: true,
      },
    });
  };

  const handleExecuteOperation = () => {
    switch (operation) {
      case 'begin':
        executeOperation('begin_transaction', {});
        break;
      case 'read':
        if (selectedTx && rowId) {
          executeOperation('read', { txId: selectedTx, rowId });
        }
        break;
      case 'write':
        if (selectedTx && rowId && writeValue) {
          try {
            const data = JSON.parse(writeValue);
            executeOperation('write', { txId: selectedTx, rowId, data });
          } catch {
            executeOperation('write', { txId: selectedTx, rowId, data: { value: writeValue } });
          }
        }
        break;
      case 'commit':
        if (selectedTx) {
          executeOperation('commit', { txId: selectedTx });
        }
        break;
      case 'abort':
        if (selectedTx) {
          executeOperation('abort', { txId: selectedTx });
        }
        break;
    }
  };

  const handleScenarioSelect = (scenarioId: string) => {
    setSelectedScenario(scenarioId);
    selectScenario(scenarioId);
  };

  return (
    <div className="mvcc-page">
      <header className="page-header">
        <Link to="/" className="back-link">
          <ArrowLeft size={20} />
          Back to Topics
        </Link>
        <div className="header-content">
          <h1>MVCC Visualization</h1>
          <p>Multi-Version Concurrency Control</p>
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
              Start Simulation
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
            <h3>Transaction Operations</h3>
            <div className="operation-grid">
              <button
                className={`op-btn ${operation === 'begin' ? 'active' : ''}`}
                onClick={() => setOperation('begin')}
              >
                <PlayCircle size={14} />
                Begin
              </button>
              <button
                className={`op-btn ${operation === 'read' ? 'active' : ''}`}
                onClick={() => setOperation('read')}
              >
                <Eye size={14} />
                Read
              </button>
              <button
                className={`op-btn ${operation === 'write' ? 'active' : ''}`}
                onClick={() => setOperation('write')}
              >
                <Edit size={14} />
                Write
              </button>
              <button
                className={`op-btn ${operation === 'commit' ? 'active' : ''}`}
                onClick={() => setOperation('commit')}
              >
                <CheckCircle size={14} />
                Commit
              </button>
              <button
                className={`op-btn ${operation === 'abort' ? 'active' : ''}`}
                onClick={() => setOperation('abort')}
              >
                <XCircle size={14} />
                Abort
              </button>
            </div>

            {operation !== 'begin' && (
              <div className="operation-params">
                <select
                  value={selectedTx}
                  onChange={(e) => setSelectedTx(e.target.value)}
                  disabled={!isConnected}
                >
                  <option value="">Select Transaction</option>
                  {activeTxs.map((tx) => (
                    <option key={tx.id} value={tx.id}>{tx.id}</option>
                  ))}
                </select>

                {(operation === 'read' || operation === 'write') && (
                  <input
                    type="text"
                    placeholder="Row ID (e.g., users:1)"
                    value={rowId}
                    onChange={(e) => setRowId(e.target.value)}
                    disabled={!isConnected}
                  />
                )}

                {operation === 'write' && (
                  <input
                    type="text"
                    placeholder='Data (e.g., {"name": "Alice"})'
                    value={writeValue}
                    onChange={(e) => setWriteValue(e.target.value)}
                    disabled={!isConnected}
                  />
                )}
              </div>
            )}

            <button
              className="execute-btn"
              onClick={handleExecuteOperation}
              disabled={!isConnected}
            >
              Execute {operation}
            </button>
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
          <div className="viz-section">
            <h3>Transactions</h3>
            <TransactionTimeline />
          </div>
          <div className="viz-section version-chains">
            <h3>Version Chains</h3>
            <VersionChainView />
          </div>
        </main>
      </div>
    </div>
  );
}

export default MVCCPage;
