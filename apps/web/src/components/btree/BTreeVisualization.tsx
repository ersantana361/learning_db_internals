import { useMemo, useCallback } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  Node,
  Edge,
  ConnectionLineType,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import { useBTreeStore, type BTreeNodeData } from '../../stores/btreeStore';
import { useSimulationStore } from '../../stores/simulationStore';
import { BTreeNode } from './BTreeNode';
import './BTreeVisualization.css';

const nodeTypes = {
  btree: BTreeNode,
};

interface LayoutNode {
  id: string;
  x: number;
  y: number;
  width: number;
}

export function BTreeVisualization() {
  const { nodes: btreeNodes, rootId } = useBTreeStore();
  const { highlights } = useSimulationStore();

  // Build highlight map for quick lookup
  const highlightMap = useMemo(() => {
    const map: Record<string, { color: string; animation?: string }> = {};
    highlights.forEach((h) => {
      if (h.type === 'node') {
        map[h.id] = { color: h.color, animation: h.animation };
      }
    });
    return map;
  }, [highlights]);

  // Calculate tree layout
  const { flowNodes, flowEdges } = useMemo(() => {
    if (!rootId || Object.keys(btreeNodes).length === 0) {
      return { flowNodes: [], flowEdges: [] };
    }

    const nodes: Node[] = [];
    const edges: Edge[] = [];
    const layoutNodes: Record<string, LayoutNode> = {};

    // BFS to calculate positions
    const nodeWidth = 100;
    const nodeHeight = 60;
    const horizontalSpacing = 40;
    const verticalSpacing = 80;

    // First pass: calculate subtree widths
    const subtreeWidths: Record<string, number> = {};

    function calculateWidth(nodeId: string): number {
      const node = btreeNodes[nodeId];
      if (!node || node.children.length === 0) {
        const width = Math.max(nodeWidth, node?.keys.length * 40 + 20);
        subtreeWidths[nodeId] = width;
        return width;
      }

      let totalWidth = 0;
      node.children.forEach((childId) => {
        totalWidth += calculateWidth(childId) + horizontalSpacing;
      });
      totalWidth -= horizontalSpacing; // Remove last spacing

      const myWidth = Math.max(nodeWidth, node.keys.length * 40 + 20);
      subtreeWidths[nodeId] = Math.max(myWidth, totalWidth);
      return subtreeWidths[nodeId];
    }

    calculateWidth(rootId);

    // Second pass: assign positions
    function assignPositions(nodeId: string, x: number, y: number): void {
      const node = btreeNodes[nodeId];
      if (!node) return;

      const width = Math.max(nodeWidth, node.keys.length * 40 + 20);
      layoutNodes[nodeId] = { id: nodeId, x, y, width };

      if (node.children.length > 0) {
        let childX = x - subtreeWidths[nodeId] / 2;

        node.children.forEach((childId) => {
          const childWidth = subtreeWidths[childId];
          assignPositions(childId, childX + childWidth / 2, y + nodeHeight + verticalSpacing);
          childX += childWidth + horizontalSpacing;
        });
      }
    }

    assignPositions(rootId, 0, 0);

    // Create React Flow nodes and edges
    Object.entries(layoutNodes).forEach(([id, layout]) => {
      const btreeNode = btreeNodes[id];
      if (!btreeNode) return;

      const highlight = highlightMap[id];

      nodes.push({
        id,
        type: 'btree',
        position: { x: layout.x - layout.width / 2, y: layout.y },
        data: {
          ...btreeNode,
          highlighted: !!highlight,
          highlightColor: highlight?.color,
          animation: highlight?.animation,
        },
      });

      // Create edges to children
      btreeNode.children.forEach((childId, index) => {
        edges.push({
          id: `${id}-${childId}`,
          source: id,
          target: childId,
          type: 'smoothstep',
          animated: highlightMap[`${id}-${childId}`] !== undefined,
          style: {
            stroke: highlightMap[`${id}-${childId}`]?.color || '#64748b',
            strokeWidth: 2,
          },
          label: String(index),
          labelStyle: { fontSize: 10, fill: '#94a3b8' },
          labelBgStyle: { fill: '#1e293b' },
        });
      });
    });

    return { flowNodes: nodes, flowEdges: edges };
  }, [btreeNodes, rootId, highlightMap]);

  const [nodes, setNodes, onNodesChange] = useNodesState(flowNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(flowEdges);

  // Update nodes when btreeNodes change
  useMemo(() => {
    setNodes(flowNodes);
    setEdges(flowEdges);
  }, [flowNodes, flowEdges, setNodes, setEdges]);

  const onInit = useCallback((instance: any) => {
    instance.fitView({ padding: 0.2 });
  }, []);

  if (!rootId || Object.keys(btreeNodes).length === 0) {
    return (
      <div className="btree-visualization empty">
        <div className="empty-state">
          <h3>No B-Tree Data</h3>
          <p>Start a simulation to visualize the B-Tree</p>
        </div>
      </div>
    );
  }

  return (
    <div className="btree-visualization">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        connectionLineType={ConnectionLineType.SmoothStep}
        fitView
        onInit={onInit}
        proOptions={{ hideAttribution: true }}
      >
        <Background color="#334155" gap={20} />
        <Controls />
        <MiniMap
          nodeColor={(node) => {
            const data = node.data as unknown as BTreeNodeData | undefined;
            if (data?.highlighted) return data.highlightColor || '#3b82f6';
            return data?.isLeaf ? '#1e3a5f' : '#1e293b';
          }}
          maskColor="rgba(0, 0, 0, 0.8)"
        />
      </ReactFlow>
    </div>
  );
}

export default BTreeVisualization;
