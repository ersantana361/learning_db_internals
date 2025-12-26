import { useCallback } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  Handle,
  Position,
} from '@xyflow/react';
import type { Node, Edge, NodeProps } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useSimulationStore } from '../../stores/simulationStore';

interface NodeVisualizerProps {
  onNodeClick?: (nodeId: string) => void;
  onCrashNode?: (nodeId: string) => void;
  onRecoverNode?: (nodeId: string) => void;
}

interface DistributedNodeData {
  label: string;
  status: string;
  role?: string;
  term?: number;
  [key: string]: unknown;
}

// Custom node component
function DistributedNode({ data }: NodeProps<Node<DistributedNodeData>>) {
  const nodeData = data as DistributedNodeData;
  const statusColor: Record<string, string> = {
    running: '#22c55e',
    crashed: '#ef4444',
    partitioned: '#f59e0b',
    byzantine: '#a855f7',
  };
  const color = statusColor[nodeData.status] || '#6b7280';

  const roleLabel = nodeData.role ? ` (${nodeData.role})` : '';

  return (
    <div
      className="distributed-node"
      style={{
        border: `3px solid ${color}`,
        borderRadius: '50%',
        width: 80,
        height: 80,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        background: nodeData.status === 'crashed' ? '#fee2e2' : '#fff',
        cursor: 'pointer',
      }}
    >
      <Handle type="target" position={Position.Top} />
      <div style={{ fontWeight: 'bold', fontSize: 14 }}>{nodeData.label}</div>
      <div style={{ fontSize: 10, color: '#666' }}>
        {nodeData.status}
        {roleLabel}
      </div>
      {nodeData.term !== undefined && (
        <div style={{ fontSize: 9, color: '#999' }}>Term: {nodeData.term}</div>
      )}
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}

const nodeTypes = {
  distributed: DistributedNode,
};

export function NodeVisualizer({
  onNodeClick,
  onCrashNode,
  onRecoverNode,
}: NodeVisualizerProps) {
  const { nodes, partitions, selectedNode, setSelectedNode } = useSimulationStore();

  // Convert nodes to React Flow nodes
  const flowNodes: Node<DistributedNodeData>[] = Object.entries(nodes).map(([id, node], index) => {
    const angle = (2 * Math.PI * index) / Object.keys(nodes).length;
    const radius = 150;
    const x = 250 + radius * Math.cos(angle);
    const y = 200 + radius * Math.sin(angle);

    return {
      id,
      type: 'distributed',
      position: { x, y },
      data: {
        label: id,
        status: node.status,
        role: node.role,
        term: node.term,
      },
      selected: id === selectedNode,
    };
  });

  // Create edges based on cluster topology (fully connected for now)
  const flowEdges: Edge[] = [];
  const nodeIds = Object.keys(nodes);

  for (let i = 0; i < nodeIds.length; i++) {
    for (let j = i + 1; j < nodeIds.length; j++) {
      const from = nodeIds[i];
      const to = nodeIds[j];

      // Check if partitioned
      const isPartitioned = partitions.some(
        (p) =>
          (p.from === from && p.to === to) ||
          (p.from === to && p.to === from)
      );

      flowEdges.push({
        id: `${from}-${to}`,
        source: from,
        target: to,
        style: {
          stroke: isPartitioned ? '#ef4444' : '#94a3b8',
          strokeWidth: 2,
          strokeDasharray: isPartitioned ? '5,5' : undefined,
        },
        animated: !isPartitioned,
      });
    }
  }

  const handleNodeClick = useCallback(
    (_event: React.MouseEvent, node: Node) => {
      setSelectedNode(node.id);
      onNodeClick?.(node.id);
    },
    [onNodeClick, setSelectedNode]
  );

  return (
    <div style={{ width: '100%', height: '400px' }}>
      <ReactFlow
        nodes={flowNodes}
        edges={flowEdges}
        nodeTypes={nodeTypes}
        onNodeClick={handleNodeClick}
        fitView
        attributionPosition="bottom-left"
      >
        <Background />
        <Controls />
      </ReactFlow>

      {selectedNode && nodes[selectedNode] && (
        <div className="node-details">
          <h4>Node: {selectedNode}</h4>
          <p>Status: {nodes[selectedNode].status}</p>
          {nodes[selectedNode].role && <p>Role: {nodes[selectedNode].role}</p>}
          {nodes[selectedNode].term !== undefined && (
            <p>Term: {nodes[selectedNode].term}</p>
          )}
          <div className="node-actions">
            {nodes[selectedNode].status === 'crashed' ? (
              <button onClick={() => onRecoverNode?.(selectedNode)}>
                Recover Node
              </button>
            ) : (
              <button onClick={() => onCrashNode?.(selectedNode)}>
                Crash Node
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
