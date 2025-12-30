import { useSimulationStore } from '../../stores/simulationStore';
import './StepTimeline.css';

export function StepTimeline() {
  const { steps, currentStep } = useSimulationStore();
  const currentIndex = currentStep?.index ?? -1;

  if (steps.length === 0) {
    return null;
  }

  return (
    <div className="step-timeline">
      <div className="timeline-track">
        {steps.map((step, index) => {
          const isActive = index === currentIndex;
          const isCompleted = index < currentIndex;

          return (
            <div
              key={index}
              className={`timeline-step ${isActive ? 'active' : ''} ${isCompleted ? 'completed' : ''}`}
            >
              <div className="step-marker">
                <span className="step-number">{index + 1}</span>
              </div>
              <div className="step-label">
                <span className="step-title-small">{step.title}</span>
              </div>
              {index < steps.length - 1 && <div className="step-connector" />}
            </div>
          );
        })}
      </div>

      <div className="progress-bar">
        <div
          className="progress-fill"
          style={{
            width: `${((currentIndex + 1) / steps.length) * 100}%`,
          }}
        />
      </div>
    </div>
  );
}
