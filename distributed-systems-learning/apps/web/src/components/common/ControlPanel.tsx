import { Play, Pause, StepForward, Square, Gauge } from 'lucide-react';
import { useSimulationStore } from '../../stores/simulationStore';

interface ControlPanelProps {
  onStart: () => void;
  onPause: () => void;
  onResume: () => void;
  onStop: () => void;
  onStep: () => void;
  onSpeedChange: (speed: number) => void;
}

export function ControlPanel({
  onStart,
  onPause,
  onResume,
  onStop,
  onStep,
  onSpeedChange,
}: ControlPanelProps) {
  const { simulation, connected } = useSimulationStore();

  return (
    <div className="control-panel">
      <div className="control-buttons">
        {!simulation.running ? (
          <button
            onClick={onStart}
            disabled={!connected}
            className="control-btn start"
            title="Start"
          >
            <Play size={20} />
          </button>
        ) : simulation.mode === 'paused' ? (
          <button
            onClick={onResume}
            className="control-btn resume"
            title="Resume"
          >
            <Play size={20} />
          </button>
        ) : (
          <button
            onClick={onPause}
            className="control-btn pause"
            title="Pause"
          >
            <Pause size={20} />
          </button>
        )}

        <button
          onClick={onStep}
          disabled={!simulation.running}
          className="control-btn step"
          title="Step Forward"
        >
          <StepForward size={20} />
        </button>

        <button
          onClick={onStop}
          disabled={!simulation.running}
          className="control-btn stop"
          title="Stop"
        >
          <Square size={20} />
        </button>
      </div>

      <div className="speed-control">
        <Gauge size={16} />
        <label>Speed:</label>
        <input
          type="range"
          min="0.1"
          max="5"
          step="0.1"
          value={simulation.speed}
          onChange={(e) => onSpeedChange(parseFloat(e.target.value))}
        />
        <span>{simulation.speed.toFixed(1)}x</span>
      </div>

      <div className="status">
        <span className={`status-indicator ${connected ? 'connected' : 'disconnected'}`}>
          {connected ? 'Connected' : 'Disconnected'}
        </span>
        <span className="mode-indicator">
          Mode: {simulation.mode}
        </span>
      </div>
    </div>
  );
}
