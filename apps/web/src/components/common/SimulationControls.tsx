import { Play, Pause, SkipForward, SkipBack, RotateCcw, Gauge } from 'lucide-react';
import { useSimulationStore } from '../../stores/simulationStore';
import './SimulationControls.css';

interface SimulationControlsProps {
  onStepForward: () => void;
  onStepBackward: () => void;
  onPlay: () => void;
  onPause: () => void;
  onReset: () => void;
  onSpeedChange: (speed: number) => void;
  disabled?: boolean;
}

export function SimulationControls({
  onStepForward,
  onStepBackward,
  onPlay,
  onPause,
  onReset,
  onSpeedChange,
  disabled = false,
}: SimulationControlsProps) {
  const { mode, speed, currentStep, connected } = useSimulationStore();
  const isDisabled = disabled || !connected;

  return (
    <div className="simulation-controls">
      <div className="step-info">
        {currentStep ? (
          <>
            <span className="step-counter">
              Step {currentStep.index + 1} of {currentStep.total}
            </span>
            <h4 className="step-title">{currentStep.title}</h4>
            <p className="step-description">{currentStep.description}</p>
          </>
        ) : (
          <p className="step-placeholder">No simulation running</p>
        )}
      </div>

      <div className="control-buttons">
        <button
          className="control-btn"
          onClick={onStepBackward}
          disabled={isDisabled || !currentStep || currentStep.index === 0}
          title="Step backward"
        >
          <SkipBack size={20} />
        </button>

        {mode === 'playing' ? (
          <button
            className="control-btn primary"
            onClick={onPause}
            disabled={isDisabled}
            title="Pause"
          >
            <Pause size={24} />
          </button>
        ) : (
          <button
            className="control-btn primary"
            onClick={onPlay}
            disabled={isDisabled || !currentStep}
            title="Play"
          >
            <Play size={24} />
          </button>
        )}

        <button
          className="control-btn"
          onClick={onStepForward}
          disabled={
            isDisabled ||
            !currentStep ||
            currentStep.index >= currentStep.total - 1
          }
          title="Step forward"
        >
          <SkipForward size={20} />
        </button>

        <button
          className="control-btn"
          onClick={onReset}
          disabled={isDisabled}
          title="Reset"
        >
          <RotateCcw size={20} />
        </button>
      </div>

      <div className="speed-control">
        <Gauge size={16} />
        <input
          type="range"
          min="0.25"
          max="4"
          step="0.25"
          value={speed}
          onChange={(e) => onSpeedChange(parseFloat(e.target.value))}
          disabled={isDisabled}
        />
        <span className="speed-value">{speed}x</span>
      </div>

      {!connected && (
        <div className="connection-status">
          <span className="status-dot disconnected" />
          Disconnected
        </div>
      )}
    </div>
  );
}
