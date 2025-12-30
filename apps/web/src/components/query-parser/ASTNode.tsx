import { Handle, Position } from '@xyflow/react';
import { motion } from 'framer-motion';
import './ASTNode.css';

interface ASTNodeProps {
  data: {
    type: string;
    value?: string;
    color: string;
    highlighted?: boolean;
    highlightColor?: string;
    animation?: string;
  };
}

export function ASTNode({ data }: ASTNodeProps) {
  const { type, value, color, highlighted, highlightColor, animation } = data;

  return (
    <motion.div
      className={`ast-node ${highlighted ? 'highlighted' : ''}`}
      style={{
        borderColor: highlighted ? highlightColor || color : color,
        boxShadow: highlighted ? `0 0 12px ${highlightColor || color}` : undefined,
      }}
      animate={animation === 'fadeIn' ? {
        opacity: [0, 1],
        scale: [0.8, 1],
        transition: { duration: 0.3 }
      } : animation === 'pulse' ? {
        scale: [1, 1.05, 1],
        transition: { duration: 0.5, repeat: 2 }
      } : undefined}
    >
      <Handle type="target" position={Position.Top} className="ast-handle" />

      <div className="node-content">
        <span className="node-type" style={{ color }}>{type}</span>
        {value && <span className="node-value">{value}</span>}
      </div>

      <Handle type="source" position={Position.Bottom} className="ast-handle" />
    </motion.div>
  );
}

export default ASTNode;
