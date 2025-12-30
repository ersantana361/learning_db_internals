import { motion } from 'framer-motion';
import { useQueryParserStore } from '../../stores/queryParserStore';
import { useSimulationStore } from '../../stores/simulationStore';
import './TokenDisplay.css';

const TOKEN_COLORS: Record<string, string> = {
  SELECT: '#3b82f6',
  FROM: '#10b981',
  WHERE: '#f59e0b',
  AND: '#8b5cf6',
  OR: '#8b5cf6',
  IDENTIFIER: '#94a3b8',
  NUMBER: '#f472b6',
  STRING: '#22d3ee',
  OPERATOR: '#ef4444',
  COMMA: '#64748b',
  STAR: '#fbbf24',
  LPAREN: '#64748b',
  RPAREN: '#64748b',
  EOF: '#475569',
};

export function TokenDisplay() {
  const { query, tokens, currentTokenIndex } = useQueryParserStore();
  const { highlights } = useSimulationStore();

  const highlightMap = new Map(
    highlights.filter(h => h.type === 'token').map(h => [h.id, h])
  );

  if (!query) {
    return (
      <div className="token-display empty">
        <p>Enter a SQL query to see tokens</p>
      </div>
    );
  }

  return (
    <div className="token-display">
      <div className="query-input">
        <span className="query-label">Query:</span>
        <code className="query-text">{query}</code>
      </div>

      <div className="tokens-container">
        <span className="tokens-label">Tokens:</span>
        <div className="tokens-list">
          {tokens.map((token, index) => {
            const highlight = highlightMap.get(`token-${index}`);
            const isHighlighted = !!highlight || index === currentTokenIndex;
            const color = TOKEN_COLORS[token.type] || '#94a3b8';

            return (
              <motion.div
                key={index}
                className={`token ${isHighlighted ? 'highlighted' : ''}`}
                style={{
                  borderColor: isHighlighted ? highlight?.color || color : 'transparent',
                  backgroundColor: `${color}20`,
                }}
                animate={isHighlighted && highlight?.animation === 'pulse' ? {
                  scale: [1, 1.1, 1],
                  transition: { duration: 0.3, repeat: 2 }
                } : undefined}
              >
                <span className="token-type" style={{ color }}>
                  {token.type}
                </span>
                {token.value && (
                  <span className="token-value">{token.value}</span>
                )}
              </motion.div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default TokenDisplay;
