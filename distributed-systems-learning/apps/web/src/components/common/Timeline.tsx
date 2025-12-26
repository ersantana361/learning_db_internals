import { useSimulationStore } from '../../stores/simulationStore';
import { Send, ArrowDown, X, Crown, CheckCircle } from 'lucide-react';

export function Timeline() {
  const { timeline } = useSimulationStore();

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'message_sent':
        return <Send size={14} className="text-blue-500" />;
      case 'message_received':
        return <ArrowDown size={14} className="text-green-500" />;
      case 'message_dropped':
        return <X size={14} className="text-red-500" />;
      case 'leader_elected':
        return <Crown size={14} className="text-yellow-500" />;
      case 'consensus_reached':
        return <CheckCircle size={14} className="text-purple-500" />;
      default:
        return null;
    }
  };

  const formatEvent = (event: typeof timeline[0]) => {
    switch (event.type) {
      case 'message_sent':
        return `${event.data.from} → ${event.data.to}: ${event.data.messageType}`;
      case 'message_received':
        return `${event.data.at} received from ${event.data.from}`;
      case 'message_dropped':
        return `Message dropped: ${event.data.reason}`;
      case 'leader_elected':
        return `${event.data.leaderId} elected leader (term ${event.data.term})`;
      case 'consensus_reached':
        return `Consensus reached on value`;
      default:
        return event.type;
    }
  };

  const recentEvents = timeline.slice(-20).reverse();

  return (
    <div className="timeline">
      <h3>Event Timeline</h3>
      <div className="timeline-events">
        {recentEvents.length === 0 ? (
          <div className="no-events">No events yet</div>
        ) : (
          recentEvents.map((event, index) => (
            <div key={index} className="timeline-event">
              <span className="event-icon">{getEventIcon(event.type)}</span>
              <span className="event-time">
                {new Date(event.time).toLocaleTimeString()}
              </span>
              <span className="event-description">{formatEvent(event)}</span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
