import { Metrics } from '../services/api';
import { Pod, Node } from '../types';
import './MetricsView.css';

interface MetricsViewProps {
  metrics: Metrics;
  resource: Pod | Node;
  resourceType: 'pod' | 'node';
}

export default function MetricsView({ metrics, resource, resourceType }: MetricsViewProps) {
  const isPod = resourceType === 'pod';
  const pod = isPod ? (resource as Pod) : null;
  const node = !isPod ? (resource as Node) : null;

  // Calculate percentages for visual bars
  let cpuPercent = 0;
  let memoryPercent = 0;
  let cpuTotal = 0;
  let memoryTotal = 0;

  if (node) {
    cpuTotal = node.cpu.total * 1000; // cores to millicores
    memoryTotal = node.memory.total * 1024; // GB to MB
    cpuPercent = (metrics.cpuUsage / cpuTotal) * 100;
    memoryPercent = (metrics.memoryUsage / memoryTotal) * 100;
  }

  const formatCPU = (millicores: number) => {
    if (millicores >= 1000) {
      return `${(millicores / 1000).toFixed(2)} cores`;
    }
    return `${millicores.toFixed(0)}m`;
  };

  const formatMemory = (mb: number) => {
    if (mb >= 1024) {
      return `${(mb / 1024).toFixed(2)} GB`;
    }
    return `${mb.toFixed(0)} MB`;
  };

  const getBarColor = (percent: number) => {
    if (percent < 50) return '#10b981'; // green
    if (percent < 80) return '#f59e0b'; // yellow/orange
    return '#ef4444'; // red
  };

  return (
    <div className="metrics-view">
      <div className="metrics-header">
        <h3>Current Resource Usage</h3>
        <div className="metrics-timestamp">
          Updated: {new Date(metrics.timestamp).toLocaleTimeString()}
        </div>
      </div>

      <div className="metrics-grid">
        {/* CPU Usage */}
        <div className="metric-card">
          <div className="metric-icon">🔥</div>
          <div className="metric-info">
            <div className="metric-label">CPU Usage</div>
            <div className="metric-value">{formatCPU(metrics.cpuUsage)}</div>
            {node && (
              <>
                <div className="metric-bar">
                  <div
                    className="metric-bar-fill"
                    style={{
                      width: `${Math.min(cpuPercent, 100)}%`,
                      backgroundColor: getBarColor(cpuPercent),
                    }}
                  />
                </div>
                <div className="metric-subtext">
                  {cpuPercent.toFixed(1)}% of {formatCPU(cpuTotal)}
                </div>
              </>
            )}
          </div>
        </div>

        {/* Memory Usage */}
        <div className="metric-card">
          <div className="metric-icon">💾</div>
          <div className="metric-info">
            <div className="metric-label">Memory Usage</div>
            <div className="metric-value">{formatMemory(metrics.memoryUsage)}</div>
            {node && (
              <>
                <div className="metric-bar">
                  <div
                    className="metric-bar-fill"
                    style={{
                      width: `${Math.min(memoryPercent, 100)}%`,
                      backgroundColor: getBarColor(memoryPercent),
                    }}
                  />
                </div>
                <div className="metric-subtext">
                  {memoryPercent.toFixed(1)}% of {formatMemory(memoryTotal)}
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Additional Info for Pods */}
      {pod && pod.containers && pod.containers.length > 0 && (
        <div className="containers-metrics">
          <h4>Containers</h4>
          <div className="containers-list-small">
            {pod.containers.map((container, idx) => (
              <div key={idx} className="container-item-small">
                <span className="container-name-small">{container.name}</span>
                <span className={`container-status-small status-${container.status.toLowerCase()}`}>
                  {container.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="metrics-note">
        <p>
          📊 Metrics are fetched from the Kubernetes metrics-server in real-time.
          {isPod && ' Pod metrics show aggregate usage across all containers.'}
        </p>
      </div>
    </div>
  );
}
