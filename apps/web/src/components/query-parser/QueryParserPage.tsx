import { useState } from 'react';
import { ArrowLeft, Code } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useWebSocket } from '../../hooks/useWebSocket';
import { SimulationControls } from '../common/SimulationControls';
import { StepTimeline } from '../common/StepTimeline';
import { TokenDisplay } from './TokenDisplay';
import { ASTVisualization } from './ASTVisualization';
import './QueryParserPage.css';

const EXAMPLE_QUERIES = [
  { name: 'Simple Select', query: 'SELECT * FROM users' },
  { name: 'With Columns', query: 'SELECT id, name, email FROM users' },
  { name: 'With Where', query: 'SELECT * FROM users WHERE age > 18' },
  { name: 'Complex', query: "SELECT id, name FROM users WHERE status = 'active' AND age >= 21" },
  { name: 'With Join', query: 'SELECT u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id' },
  { name: 'With Order', query: 'SELECT * FROM products ORDER BY price LIMIT 10' },
];

export function QueryParserPage() {
  const [queryInput, setQueryInput] = useState('SELECT * FROM users WHERE id = 1');
  const [selectedExample, setSelectedExample] = useState('');

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
  } = useWebSocket();

  const handleParseQuery = () => {
    startSimulation({
      project: 'query-parser',
      parameters: {
        query: queryInput,
      },
    });
  };

  const handleSelectExample = (query: string, name: string) => {
    setQueryInput(query);
    setSelectedExample(name);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && e.ctrlKey) {
      handleParseQuery();
    }
  };

  return (
    <div className="query-parser-page">
      <header className="page-header">
        <Link to="/" className="back-link">
          <ArrowLeft size={20} />
          Back to Topics
        </Link>
        <div className="header-content">
          <h1>Query Parser Visualization</h1>
          <p>SQL Lexing and Parsing</p>
        </div>
        <div className="connection-badge">
          <span className={`status-dot ${isConnected ? 'connected' : 'disconnected'}`} />
          {isConnecting ? 'Connecting...' : isConnected ? 'Connected' : 'Disconnected'}
        </div>
      </header>

      <div className="page-content">
        <aside className="control-panel">
          <section className="panel-section">
            <h3>Query Input</h3>
            <textarea
              className="query-textarea"
              value={queryInput}
              onChange={(e) => setQueryInput(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Enter SQL query..."
              disabled={!isConnected}
              rows={4}
            />
            <button
              className="parse-btn"
              onClick={handleParseQuery}
              disabled={!isConnected || !queryInput.trim()}
            >
              <Code size={16} />
              Parse Query
            </button>
            <span className="hint">Ctrl+Enter to parse</span>
          </section>

          <section className="panel-section">
            <h3>Example Queries</h3>
            <div className="example-list">
              {EXAMPLE_QUERIES.map((example) => (
                <button
                  key={example.name}
                  className={`example-btn ${selectedExample === example.name ? 'selected' : ''}`}
                  onClick={() => handleSelectExample(example.query, example.name)}
                  disabled={!isConnected}
                >
                  <span className="example-name">{example.name}</span>
                  <code className="example-query">{example.query}</code>
                </button>
              ))}
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
          <div className="viz-section tokens">
            <h3>Tokens</h3>
            <TokenDisplay />
          </div>
          <div className="viz-section ast">
            <h3>Abstract Syntax Tree</h3>
            <ASTVisualization />
          </div>
        </main>
      </div>
    </div>
  );
}

export default QueryParserPage;
