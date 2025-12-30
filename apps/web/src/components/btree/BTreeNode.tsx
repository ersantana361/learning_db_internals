import { Handle, Position } from '@xyflow/react';
import { motion } from 'framer-motion';
import type { BTreeNodeData } from '../../stores/btreeStore';
import './BTreeNode.css';

interface BTreeNodeProps {
  data: BTreeNodeData & {
    highlighted?: boolean;
    highlightColor?: string;
    animation?: string;
  };
}

export function BTreeNode({ data }: BTreeNodeProps) {
  const { keys, isLeaf, highlighted, highlightColor, animation } = data;

  const nodeClass = `btree-node ${isLeaf ? 'leaf' : 'internal'} ${highlighted ? 'highlighted' : ''}`;

  const animationVariants = {
    pulse: {
      scale: [1, 1.05, 1],
      transition: { duration: 0.5, repeat: 2 },
    },
    shake: {
      x: [0, -5, 5, -5, 5, 0],
      transition: { duration: 0.4 },
    },
    fadeIn: {
      opacity: [0, 1],
      scale: [0.8, 1],
      transition: { duration: 0.3 },
    },
  };

  return (
    <motion.div
      className={nodeClass}
      style={highlighted && highlightColor ? { borderColor: highlightColor, boxShadow: `0 0 12px ${highlightColor}` } : undefined}
      variants={animationVariants}
      animate={animation || undefined}
    >
      <Handle type="target" position={Position.Top} className="btree-handle" />

      <div className="node-keys">
        {keys.map((key, index) => (
          <div key={index} className="key-cell">
            <span className="key-value">{key}</span>
          </div>
        ))}
      </div>

      {!isLeaf && (
        <Handle type="source" position={Position.Bottom} className="btree-handle" />
      )}
    </motion.div>
  );
}

export default BTreeNode;
