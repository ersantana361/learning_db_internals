import { useMemo, useCallback } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  useNodesState,
  useEdgesState,
  Node,
  Edge,
  ConnectionLineType,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import { useQueryParserStore } from '../../stores/queryParserStore';
import { useSimulationStore } from '../../stores/simulationStore';
import { ASTNode as ASTNodeComponent } from './ASTNode';
import './ASTVisualization.css';

const nodeTypes = {
  ast: ASTNodeComponent,
};

const NODE_COLORS: Record<string, string> = {
  STATEMENT: '#3b82f6',
  SELECT: '#3b82f6',
  FROM: '#10b981',
  WHERE: '#f59e0b',
  COLUMNS: '#8b5cf6',
  COLUMN: '#a78bfa',
  TABLE: '#14b8a6',
  CONDITION: '#f97316',
  BINARY_EXPR: '#ef4444',
  LITERAL: '#f472b6',
  IDENTIFIER: '#94a3b8',
  ORDER_BY: '#06b6d4',
  LIMIT: '#84cc16',
  JOIN: '#10b981',
};

export function ASTVisualization() {
  const { astNodes, astRoot } = useQueryParserStore();
  const { highlights } = useSimulationStore();

  const highlightMap = useMemo(() => {
    const map: Record<string, { color: string; animation?: string }> = {};
    highlights.forEach((h) => {
      if (h.type === 'node') {
        map[h.id] = { color: h.color, animation: h.animation };
      }
    });
    return map;
  }, [highlights]);

  const { flowNodes, flowEdges } = useMemo(() => {
    if (!astRoot || Object.keys(astNodes).length === 0) {
      return { flowNodes: [], flowEdges: [] };
    }

    const nodes: Node[] = [];
    const edges: Edge[] = [];
    const positions: Record<string, { x: number; y: number }> = {};

    // Calculate subtree widths
    const subtreeWidths: Record<string, number> = {};
    const nodeWidth = 120;
    const nodeHeight = 50;
    const horizontalSpacing = 30;
    const verticalSpacing = 80;

    function calculateWidth(nodeId: string): number {
      const node = astNodes[nodeId];
      if (!node || node.children.length === 0) {
        subtreeWidths[nodeId] = nodeWidth;
        return nodeWidth;
      }

      let totalWidth = 0;
      node.children.forEach((childId) => {
        totalWidth += calculateWidth(childId) + horizontalSpacing;
      });
      totalWidth -= horizontalSpacing;

      subtreeWidths[nodeId] = Math.max(nodeWidth, totalWidth);
      return subtreeWidths[nodeId];
    }

    calculateWidth(astRoot);

    // Assign positions
    function assignPositions(nodeId: string, x: number, y: number): void {
      const node = astNodes[nodeId];
      if (!node) return;

      positions[nodeId] = { x: x - nodeWidth / 2, y };

      if (node.children.length > 0) {
        let childX = x - subtreeWidths[nodeId] / 2;

        node.children.forEach((childId) => {
          const childWidth = subtreeWidths[childId];
          assignPositions(childId, childX + childWidth / 2, y + nodeHeight + verticalSpacing);
          childX += childWidth + horizontalSpacing;
        });
      }
    }

    assignPositions(astRoot, 0, 0);

    // Create React Flow nodes and edges
    Object.entries(positions).forEach(([id, pos]) => {
      const astNode = astNodes[id];
      if (!astNode) return;

      const highlight = highlightMap[id];
      const baseColor = NODE_COLORS[astNode.type] || '#64748b';

      nodes.push({
        id,
        type: 'ast',
        position: pos,
        data: {
          type: astNode.type,
          value: astNode.value,
          color: baseColor,
          highlighted: !!highlight,
          highlightColor: highlight?.color,
          animation: highlight?.animation,
        },
      });

      astNode.children.forEach((childId) => {
        edges.push({
          id: `${id}-${childId}`,
          source: id,
          target: childId,
          type: 'smoothstep',
          style: {
            stroke: '#64748b',
            strokeWidth: 2,
          },
        });
      });
    });

    return { flowNodes: nodes, flowEdges: edges };
  }, [astNodes, astRoot, highlightMap]);

  const [nodes, setNodes, onNodesChange] = useNodesState(flowNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(flowEdges);

  useMemo(() => {
    setNodes(flowNodes);
    setEdges(flowEdges);
  }, [flowNodes, flowEdges, setNodes, setEdges]);

  const onInit = useCallback((instance: { fitView: (options: { padding: number }) => void }) => {
    instance.fitView({ padding: 0.2 });
  }, []);

  if (!astRoot || Object.keys(astNodes).length === 0) {
    return (
      <div className="ast-visualization empty">
        <div className="empty-state">
          <h3>No AST</h3>
          <p>Parse a query to see the Abstract Syntax Tree</p>
        </div>
      </div>
    );
  }

  return (
    <div className="ast-visualization">
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
      </ReactFlow>
    </div>
  );
}

export default ASTVisualization;
