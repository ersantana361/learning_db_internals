import { useEffect, useRef, useCallback, useState } from 'react';
import { useSimulationStore } from '../stores/simulationStore';

interface UseWebSocketOptions {
  url: string;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (error: Event) => void;
}

export function useWebSocket({ url, onOpen, onClose, onError }: UseWebSocketOptions) {
  const ws = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    setNodes,
    setSimulationState,
    addTimelineEvent,
    setConnected,
  } = useSimulationStore();

  useEffect(() => {
    const connect = () => {
      try {
        ws.current = new WebSocket(url);

        ws.current.onopen = () => {
          setIsConnected(true);
          setConnected(true);
          setError(null);
          onOpen?.();
        };

        ws.current.onclose = () => {
          setIsConnected(false);
          setConnected(false);
          onClose?.();
          // Attempt to reconnect after 3 seconds
          setTimeout(connect, 3000);
        };

        ws.current.onerror = (event) => {
          setError('WebSocket connection error');
          onError?.(event);
        };

        ws.current.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data);
            handleMessage(msg);
          } catch (e) {
            console.error('Failed to parse message:', e);
          }
        };
      } catch (e) {
        setError('Failed to connect');
      }
    };

    connect();

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [url]);

  const handleMessage = (msg: any) => {
    switch (msg.type) {
      case 'simulation_state':
        setSimulationState({
          mode: msg.mode,
          speed: msg.speed,
          virtualTime: msg.virtualTime,
          running: msg.running,
        });
        if (msg.nodes) {
          setNodes(msg.nodes);
        }
        break;

      case 'node_state_update':
        // Update specific node
        break;

      case 'message_sent':
        addTimelineEvent({
          type: 'message_sent',
          time: msg.time || Date.now(),
          data: {
            from: msg.from,
            to: msg.to,
            messageType: msg.messageType,
            messageId: msg.messageId,
          },
        });
        break;

      case 'message_received':
        addTimelineEvent({
          type: 'message_received',
          time: msg.time || Date.now(),
          data: {
            at: msg.at,
            from: msg.from,
            messageId: msg.messageId,
            latency: msg.latency,
          },
        });
        break;

      case 'message_dropped':
        addTimelineEvent({
          type: 'message_dropped',
          time: msg.time || Date.now(),
          data: {
            messageId: msg.messageId,
            reason: msg.reason,
          },
        });
        break;

      case 'leader_elected':
        addTimelineEvent({
          type: 'leader_elected',
          time: msg.time || Date.now(),
          data: {
            leaderId: msg.leaderId,
            term: msg.term,
          },
        });
        break;

      case 'error':
        setError(msg.message);
        break;

      default:
        console.log('Unknown message type:', msg.type);
    }
  };

  const send = useCallback((msg: object) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(msg));
    }
  }, []);

  // Simulation control functions
  const startSimulation = useCallback((project: string, scenario?: string, config?: any) => {
    send({
      type: 'start_simulation',
      project,
      scenario,
      config,
    });
  }, [send]);

  const pauseSimulation = useCallback(() => {
    send({ type: 'pause_simulation' });
  }, [send]);

  const resumeSimulation = useCallback(() => {
    send({ type: 'resume_simulation' });
  }, [send]);

  const stopSimulation = useCallback(() => {
    send({ type: 'stop_simulation' });
  }, [send]);

  const stepForward = useCallback(() => {
    send({ type: 'step_forward' });
  }, [send]);

  const setSpeed = useCallback((speed: number) => {
    send({ type: 'set_speed', speed });
  }, [send]);

  const injectCrash = useCallback((nodeId: string) => {
    send({ type: 'inject_crash', nodeId });
  }, [send]);

  const recoverNode = useCallback((nodeId: string) => {
    send({ type: 'recover_node', nodeId });
  }, [send]);

  const injectPartition = useCallback((from: string, to: string, bidirectional = false) => {
    send({ type: 'inject_partition', from, to, bidirectional });
  }, [send]);

  const healPartition = useCallback((from: string, to: string, bidirectional = false) => {
    send({ type: 'heal_partition', from, to, bidirectional });
  }, [send]);

  const getState = useCallback(() => {
    send({ type: 'get_state' });
  }, [send]);

  return {
    isConnected,
    error,
    send,
    startSimulation,
    pauseSimulation,
    resumeSimulation,
    stopSimulation,
    stepForward,
    setSpeed,
    injectCrash,
    recoverNode,
    injectPartition,
    healPartition,
    getState,
  };
}
