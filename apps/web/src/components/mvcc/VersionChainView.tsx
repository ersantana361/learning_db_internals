import { motion } from 'framer-motion';
import { ArrowRight } from 'lucide-react';
import { useMVCCStore } from '../../stores/mvccStore';
import { useSimulationStore } from '../../stores/simulationStore';
import './VersionChainView.css';

export function VersionChainView() {
  const { rows, versions } = useMVCCStore();
  const { highlights } = useSimulationStore();

  const highlightMap = new Map(
    highlights.filter(h => h.type === 'cell').map(h => [h.id, h])
  );

  const rowList = Object.values(rows);

  if (rowList.length === 0) {
    return (
      <div className="version-chain-view empty">
        <p>No rows in database</p>
      </div>
    );
  }

  return (
    <div className="version-chain-view">
      {rowList.map((row) => (
        <div key={row.id} className="row-chain">
          <div className="row-header">
            <span className="row-id">{row.id}</span>
            <span className="version-count">{row.versionChain.length} version(s)</span>
          </div>
          <div className="chain-container">
            {row.versionChain.map((verID, index) => {
              const version = versions[verID];
              if (!version) return null;

              const highlight = highlightMap.get(verID);
              const isHighlighted = !!highlight;

              return (
                <div key={verID} className="version-wrapper">
                  {index > 0 && (
                    <div className="chain-arrow">
                      <ArrowRight size={16} />
                    </div>
                  )}
                  <motion.div
                    className={`version-card ${isHighlighted ? 'highlighted' : ''} ${version.deletedBy ? 'deleted' : ''}`}
                    style={isHighlighted && highlight ? {
                      borderColor: highlight.color,
                      boxShadow: `0 0 12px ${highlight.color}`
                    } : undefined}
                    animate={isHighlighted && highlight?.animation === 'pulse' ? {
                      scale: [1, 1.05, 1],
                      transition: { duration: 0.5, repeat: 2 }
                    } : isHighlighted && highlight?.animation === 'fadeIn' ? {
                      opacity: [0, 1],
                      scale: [0.9, 1],
                      transition: { duration: 0.3 }
                    } : undefined}
                  >
                    <div className="version-header">
                      <span className="version-id">{verID}</span>
                      {index === 0 && <span className="current-badge">Current</span>}
                    </div>
                    <div className="version-meta">
                      <span>Created by: {version.createdBy}</span>
                      <span>At: {version.createdAt}</span>
                    </div>
                    <div className="version-data">
                      {Object.entries(version.data).map(([key, value]) => (
                        <div key={key} className="data-field">
                          <span className="field-key">{key}:</span>
                          <span className="field-value">{JSON.stringify(value)}</span>
                        </div>
                      ))}
                    </div>
                    {version.deletedBy && (
                      <div className="deleted-info">
                        Deleted by {version.deletedBy} at {version.deletedAt}
                      </div>
                    )}
                  </motion.div>
                </div>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}

export default VersionChainView;
