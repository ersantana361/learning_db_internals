import { useEffect, useRef, useCallback, useState } from 'react';
import { useSimulationStore } from '../stores/simulationStore';
import { useBTreeStore } from '../stores/btreeStore';
import { useMVCCStore } from '../stores/mvccStore';
import { useQueryParserStore } from '../stores/queryParserStore';
import type { Message, SimulationConfig, StepInfo } from '../types';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
const RECONNECT_DELAY = 3000;

interface UseWebSocketOptions {
  autoConnect?: boolean;
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const { autoConnect = true } = options;
  const ws = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [isConnecting, setIsConnecting] = useState(false);

  // Get store actions
  const {
    setConnected,
    setProject,
    setMode,
    setSpeed,
    setCurrentStep,
    setTotalSteps,
    setSteps,
    setHighlights,
  } = useSimulationStore();

  const btreeStore = useBTreeStore();
  const mvccStore = useMVCCStore();
  const parserStore = useQueryParserStore();

  // Handle incoming messages
  const handleMessage = useCallback(
    (event: MessageEvent) => {
      try {
        const msg: Message = JSON.parse(event.data);

        switch (msg.type) {
          case 'simulation_state': {
            const payload = msg.payload as {
              state: {
                project: string;
                mode: 'idle' | 'playing' | 'paused' | 'step';
                speed: number;
                currentStep: StepInfo | null;
                totalSteps: number;
                data: Record<string, unknown>;
              };
              steps: StepInfo[];
            };
            setProject(payload.state.project);
            setMode(payload.state.mode);
            setSpeed(payload.state.speed);
            setCurrentStep(payload.state.currentStep);
            setTotalSteps(payload.state.totalSteps);
            setSteps(payload.steps || []);

            // Route data to appropriate store
            if (payload.state.project === 'btree' && payload.state.data) {
              const data = payload.state.data as {
                nodes?: Record<string, unknown>;
                rootId?: string;
                order?: number;
              };
              if (data.nodes) btreeStore.setNodes(data.nodes as never);
              if (data.rootId !== undefined) btreeStore.setRoot(data.rootId);
              if (data.order) btreeStore.setOrder(data.order);
            } else if (payload.state.project === 'mvcc' && payload.state.data) {
              const data = payload.state.data as {
                transactions?: Record<string, unknown>;
                versions?: Record<string, unknown>;
                rows?: Record<string, unknown>;
              };
              if (data.transactions)
                mvccStore.setTransactions(data.transactions as never);
              if (data.versions) mvccStore.setVersions(data.versions as never);
              if (data.rows) mvccStore.setRows(data.rows as never);
            } else if (
              payload.state.project === 'query-parser' &&
              payload.state.data
            ) {
              const data = payload.state.data as {
                query?: string;
                tokens?: unknown[];
                astNodes?: Record<string, unknown>;
                astRoot?: string;
              };
              if (data.query) parserStore.setQuery(data.query);
              if (data.tokens) parserStore.setTokens(data.tokens as never);
              if (data.astNodes)
                parserStore.setAST(data.astNodes as never, data.astRoot || null);
            }
            break;
          }

          case 'step_update': {
            const payload = msg.payload as {
              project: string;
              step: StepInfo;
              highlights: Array<{
                type: string;
                id: string;
                color: string;
                animation?: string;
                duration?: number;
              }>;
              data: Record<string, unknown>;
            };
            setCurrentStep(payload.step);
            setHighlights(payload.highlights as never);

            // Route data to appropriate store
            if (payload.project === 'btree' && payload.data) {
              const data = payload.data as {
                nodes?: Record<string, unknown>;
                rootId?: string;
                path?: string[];
                stats?: unknown;
              };
              if (data.nodes) btreeStore.setNodes(data.nodes as never);
              if (data.rootId !== undefined) btreeStore.setRoot(data.rootId);
              if (data.path) btreeStore.setCurrentPath(data.path);
              if (data.stats) btreeStore.setStats(data.stats as never);
            } else if (payload.project === 'mvcc' && payload.data) {
              const data = payload.data as {
                transactions?: Record<string, unknown>;
                versions?: Record<string, unknown>;
                rows?: Record<string, unknown>;
                snapshot?: number;
              };
              if (data.transactions)
                mvccStore.setTransactions(data.transactions as never);
              if (data.versions) mvccStore.setVersions(data.versions as never);
              if (data.rows) mvccStore.setRows(data.rows as never);
              if (data.snapshot !== undefined)
                mvccStore.setSnapshot(data.snapshot);
            } else if (payload.project === 'query-parser' && payload.data) {
              const data = payload.data as {
                tokens?: unknown[];
                currentTokenIndex?: number;
                astNodes?: Record<string, unknown>;
                astRoot?: string;
                parsePhase?: string;
              };
              if (data.tokens) parserStore.setTokens(data.tokens as never);
              if (data.currentTokenIndex !== undefined)
                parserStore.highlightToken(data.currentTokenIndex);
              if (data.astNodes)
                parserStore.setAST(data.astNodes as never, data.astRoot || null);
              if (data.parsePhase)
                parserStore.setParsePhase(data.parsePhase as never);
            }
            break;
          }

          case 'error': {
            const payload = msg.payload as { code: string; message: string };
            console.error(`WebSocket error: ${payload.code} - ${payload.message}`);
            break;
          }

          default:
            console.log('Unknown message type:', msg.type);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    },
    [
      setProject,
      setMode,
      setSpeed,
      setCurrentStep,
      setTotalSteps,
      setSteps,
      setHighlights,
      btreeStore,
      mvccStore,
      parserStore,
    ]
  );

  // Connect to WebSocket
  const connect = useCallback(() => {
    if (ws.current?.readyState === WebSocket.OPEN || isConnecting) {
      return;
    }

    setIsConnecting(true);

    try {
      ws.current = new WebSocket(WS_URL);

      ws.current.onopen = () => {
        console.log('WebSocket connected');
        setConnected(true);
        setIsConnecting(false);
      };

      ws.current.onclose = () => {
        console.log('WebSocket disconnected');
        setConnected(false);
        setIsConnecting(false);

        // Attempt reconnection
        if (autoConnect) {
          reconnectTimeout.current = setTimeout(() => {
            connect();
          }, RECONNECT_DELAY);
        }
      };

      ws.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setIsConnecting(false);
      };

      ws.current.onmessage = handleMessage;
    } catch (error) {
      console.error('Error creating WebSocket:', error);
      setIsConnecting(false);
    }
  }, [autoConnect, handleMessage, setConnected, isConnecting]);

  // Disconnect
  const disconnect = useCallback(() => {
    if (reconnectTimeout.current) {
      clearTimeout(reconnectTimeout.current);
    }
    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }
    setConnected(false);
  }, [setConnected]);

  // Send message helper
  const send = useCallback((type: string, payload?: unknown) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      const message: Message = { type: type as never, payload: payload as never };
      ws.current.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket not connected');
    }
  }, []);

  // Control functions
  const startSimulation = useCallback(
    (config: SimulationConfig) => {
      send('start_simulation', { config });
    },
    [send]
  );

  const stepForward = useCallback(() => send('step_forward'), [send]);
  const stepBackward = useCallback(() => send('step_backward'), [send]);
  const play = useCallback(() => send('resume_simulation'), [send]);
  const pause = useCallback(() => send('pause_simulation'), [send]);
  const reset = useCallback(() => send('reset'), [send]);

  const setPlaybackSpeed = useCallback(
    (speed: number) => {
      send('set_speed', { speed });
    },
    [send]
  );

  const executeOperation = useCallback(
    (operation: string, params: Record<string, unknown>) => {
      send('execute_operation', { operation, params });
    },
    [send]
  );

  const selectScenario = useCallback(
    (scenarioId: string) => {
      send('select_scenario', { scenarioId });
    },
    [send]
  );

  const getState = useCallback(() => send('get_state'), [send]);

  // Auto-connect on mount
  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [autoConnect, connect, disconnect]);

  return {
    isConnected: useSimulationStore((s) => s.connected),
    isConnecting,
    connect,
    disconnect,
    startSimulation,
    stepForward,
    stepBackward,
    play,
    pause,
    reset,
    setSpeed: setPlaybackSpeed,
    executeOperation,
    selectScenario,
    getState,
  };
}
