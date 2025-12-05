import { Node, Pod } from '../types';
import './DetailPanel.css';

interface DetailPanelProps {
  resource: Pod | Node | null;
  resourceType: 'pod' | 'node' | null;
  onClose: () => void;
}

export default function DetailPanel({ resource, resourceType, onClose }: DetailPanelProps) {
  if (!resource || !resourceType) return null;

  const isPod = resourceType === 'pod';
  const pod = isPod ? (resource as Pod) : null;
  const node = !isPod ? (resource as Node) : null;

  return (
    <>
      {/* Backdrop */}
      <div className="detail-panel-backdrop" onClick={onClose} />

      {/* Panel */}
      <div className="detail-panel">
        {/* Header */}
        <div className="detail-panel-header">
          <h2>{isPod ? '📦 Pod Details' : '🖥️ Node Details'}</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        {/* Content */}
        <div className="detail-panel-content">
          {isPod && pod ? (
            // Pod Details
            <>
              <section>
                <h3>Basic Info</h3>
                <div className="info-grid">
                  <div className="info-item">
                    <span className="label">Name:</span>
                    <span className="value">{pod.name}</span>
                  </div>
                  <div className="info-item">
                    <span className="label">Namespace:</span>
                    <span className="value namespace-badge">{pod.namespace}</span>
                  </div>
                  <div className="info-item">
                    <span className="label">Status:</span>
                    <span className={`value status-badge status-${pod.status.toLowerCase()}`}>
                      {pod.status}
                    </span>
                  </div>
                  <div className="info-item">
                    <span className="label">Node:</span>
                    <span className="value">{pod.nodeName || 'Not assigned'}</span>
                  </div>
                  <div className="info-item">
                    <span className="label">Created:</span>
                    <span className="value">{new Date(pod.createdAt).toLocaleString()}</span>
                  </div>
                </div>
              </section>

              {pod.containers && pod.containers.length > 0 && (
                <section>
                  <h3>Containers ({pod.containers.length})</h3>
                  <div className="containers-list">
                    {pod.containers.map((container, idx) => (
                      <div key={idx} className="container-item">
                        <div className="container-name">{container.name}</div>
                        <div className="container-details">
                          <span className={`status-badge status-${container.status.toLowerCase()}`}>
                            {container.status}
                          </span>
                          <span className="restart-count">
                            Restarts: {container.restarts}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                </section>
              )}

              <section className="actions-section">
                <h3>Actions</h3>
                <div className="action-buttons">
                  <button className="action-button" disabled title="Coming soon">
                    📋 Describe
                  </button>
                  <button className="action-button" disabled title="Coming soon">
                    📜 View Logs
                  </button>
                  <button className="action-button" disabled title="Coming soon">
                    📊 Metrics
                  </button>
                </div>
              </section>
            </>
          ) : node ? (
            // Node Details
            <>
              <section>
                <h3>Basic Info</h3>
                <div className="info-grid">
                  <div className="info-item">
                    <span className="label">Name:</span>
                    <span className="value">{node.name}</span>
                  </div>
                  <div className="info-item">
                    <span className="label">Status:</span>
                    <span className={`value status-badge status-${node.status.toLowerCase()}`}>
                      {node.status}
                    </span>
                  </div>
                  <div className="info-item">
                    <span className="label">Pods:</span>
                    <span className="value">{node.pods?.length || 0}</span>
                  </div>
                </div>
              </section>

              <section>
                <h3>Resources</h3>
                <div className="resource-stats">
                  <div className="resource-item">
                    <div className="resource-label">CPU</div>
                    <div className="resource-bar">
                      <div
                        className="resource-fill cpu"
                        style={{ width: `${(node.cpu.used / node.cpu.total) * 100}%` }}
                      />
                    </div>
                    <div className="resource-text">
                      {node.cpu.used.toFixed(2)} / {node.cpu.total.toFixed(2)} cores
                    </div>
                  </div>
                  <div className="resource-item">
                    <div className="resource-label">Memory</div>
                    <div className="resource-bar">
                      <div
                        className="resource-fill memory"
                        style={{ width: `${(node.memory.used / node.memory.total) * 100}%` }}
                      />
                    </div>
                    <div className="resource-text">
                      {node.memory.used.toFixed(2)} / {node.memory.total.toFixed(2)} GB
                    </div>
                  </div>
                </div>
              </section>

              {node.labels && Object.keys(node.labels).length > 0 && (
                <section>
                  <h3>Labels</h3>
                  <div className="labels-list">
                    {Object.entries(node.labels).map(([key, value]) => (
                      <div key={key} className="label-item">
                        <span className="label-key">{key}:</span>
                        <span className="label-value">{value}</span>
                      </div>
                    ))}
                  </div>
                </section>
              )}

              <section className="actions-section">
                <h3>Actions</h3>
                <div className="action-buttons">
                  <button className="action-button" disabled title="Coming soon">
                    📋 Describe
                  </button>
                  <button className="action-button" disabled title="Coming soon">
                    📊 Metrics
                  </button>
                  <button className="action-button" disabled title="Coming soon">
                    🔧 Cordon/Drain
                  </button>
                </div>
              </section>
            </>
          ) : null}
        </div>
      </div>
    </>
  );
}
