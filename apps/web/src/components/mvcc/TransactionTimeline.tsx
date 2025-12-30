import { motion } from 'framer-motion';
import { useMVCCStore } from '../../stores/mvccStore';
import { useSimulationStore } from '../../stores/simulationStore';
import './TransactionTimeline.css';

export function TransactionTimeline() {
  const { transactions } = useMVCCStore();
  const { highlights } = useSimulationStore();

  const highlightMap = new Map(
    highlights.filter(h => h.type === 'row').map(h => [h.id, h])
  );

  const txList = Object.values(transactions).sort((a, b) => a.startTime - b.startTime);

  if (txList.length === 0) {
    return (
      <div className="transaction-timeline empty">
        <p>No transactions</p>
      </div>
    );
  }

  return (
    <div className="transaction-timeline">
      <div className="timeline-header">
        <span>Transaction</span>
        <span>Status</span>
        <span>Start</span>
        <span>Commit</span>
      </div>
      <div className="timeline-body">
        {txList.map((tx) => {
          const highlight = highlightMap.get(tx.id);
          const isHighlighted = !!highlight;

          return (
            <motion.div
              key={tx.id}
              className={`timeline-row ${tx.status} ${isHighlighted ? 'highlighted' : ''}`}
              style={isHighlighted && highlight ? { borderColor: highlight.color } : undefined}
              animate={isHighlighted && highlight?.animation === 'pulse' ? {
                scale: [1, 1.02, 1],
                transition: { duration: 0.5, repeat: 2 }
              } : undefined}
            >
              <span className="tx-id">{tx.id}</span>
              <span className={`tx-status ${tx.status}`}>{tx.status}</span>
              <span className="tx-time">{tx.startTime}</span>
              <span className="tx-time">{tx.commitTime ?? '-'}</span>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}

export default TransactionTimeline;
