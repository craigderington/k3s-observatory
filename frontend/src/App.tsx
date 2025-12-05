import { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import Scene from './components/Scene';
import ToastContainer, { Toast } from './components/ToastContainer';
import DetailPanel from './components/DetailPanel';
import { fetchNodes, fetchPods } from './services/api';
import { useWebSocket } from './hooks/useWebSocket';
import { Node, Pod } from './types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCircle, faServer, faCube, faCircleDot, faLayerGroup } from '@fortawesome/free-solid-svg-icons';
import './App.css';

export default function App() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [pods, setPods] = useState<Pod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const [selectedNamespace, setSelectedNamespace] = useState<string>('all');
  const [hasConnectedOnce, setHasConnectedOnce] = useState(false);
  const [selectedResource, setSelectedResource] = useState<Pod | Node | null>(null);
  const [selectedResourceType, setSelectedResourceType] = useState<'pod' | 'node' | null>(null);

  // Track raw pods before redistribution
  const [rawPods, setRawPods] = useState<Pod[]>([]);

  // Add toast notification
  const addToast = useCallback((message: string, type: Toast['type'] = 'info') => {
    const toast: Toast = {
      id: `${Date.now()}-${Math.random()}`,
      message,
      type,
      duration: 5000,
    };
    setToasts((prev) => [...prev, toast]);
  }, []);

  // Remove toast
  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  // Handle pod click
  const handlePodClick = useCallback((pod: Pod) => {
    setSelectedResource(pod);
    setSelectedResourceType('pod');
  }, []);

  // Handle node click
  const handleNodeClick = useCallback((node: Node) => {
    setSelectedResource(node);
    setSelectedResourceType('node');
  }, []);

  // Handle detail panel close
  const handleDetailPanelClose = useCallback(() => {
    setSelectedResource(null);
    setSelectedResourceType(null);
  }, []);

  // Helper function to redistribute pods evenly around their nodes
  const redistributePodsHelper = useCallback((pods: Pod[], nodes: Node[]): Pod[] => {
    // Create a map of node positions
    const nodePositions = new Map(nodes.map(n => [n.name, n.position]));

    // Group pods by node
    const podsByNode = new Map<string, Pod[]>();
    pods.forEach(pod => {
      if (pod.nodeName) {
        const nodePods = podsByNode.get(pod.nodeName) || [];
        nodePods.push(pod);
        podsByNode.set(pod.nodeName, nodePods);
      }
    });

    // Recalculate positions for pods on each node
    return pods.map(pod => {
      if (!pod.nodeName) return pod; // Unscheduled pods stay at origin

      const nodePos = nodePositions.get(pod.nodeName) || { x: 0, y: 0, z: 0 };
      const nodePods = podsByNode.get(pod.nodeName) || [];
      const podIndex = nodePods.findIndex(p => p.id === pod.id);
      const totalPods = nodePods.length;

      if (totalPods === 0) return pod;

      const angle = (podIndex * 2.0 * Math.PI) / totalPods;
      const orbitRadius = 3.0;

      return {
        ...pod,
        position: {
          x: nodePos.x + orbitRadius * Math.cos(angle),
          y: nodePos.y,
          z: nodePos.z + orbitRadius * Math.sin(angle),
        },
      };
    });
  }, []);

  // Redistribute pods whenever rawPods or nodes change
  useEffect(() => {
    if (rawPods.length === 0) {
      setPods([]);
      return;
    }

    if (nodes.length === 0) {
      setPods(rawPods); // Show pods even if nodes aren't loaded yet
      return;
    }

    const redistributed = redistributePodsHelper(rawPods, nodes);
    setPods(redistributed);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rawPods, nodes]);

  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((event: { type: string; data: any }) => {
    const { type, data } = event;

    switch (type) {
      case 'pod_added':
        if (data.pod) {
          setRawPods((prev) => {
            // Check if pod already exists
            if (prev.find((p) => p.id === data.pod.id)) {
              return prev;
            }
            addToast(`Pod created: ${data.pod.namespace}/${data.pod.name}`, 'success');
            return [...prev, data.pod];
          });
        }
        break;

      case 'pod_modified':
        if (data.pod) {
          setRawPods((prev) =>
            prev.map((p) => (p.id === data.pod.id ? data.pod : p))
          );
        }
        break;

      case 'pod_deleted':
        if (data.pod) {
          setRawPods((prev) => {
            const pod = prev.find((p) => p.id === data.pod.id);
            if (pod) {
              addToast(`Pod deleted: ${pod.namespace}/${pod.name}`, 'warning');
            }
            return prev.filter((p) => p.id !== data.pod.id);
          });
        }
        break;

      case 'node_added':
        if (data.node) {
          setNodes((prev) => {
            if (prev.find((n) => n.id === data.node.id)) {
              return prev;
            }
            addToast(`Node added: ${data.node.name}`, 'success');
            return [...prev, data.node];
          });
        }
        break;

      case 'node_modified':
        if (data.node) {
          setNodes((prev) =>
            prev.map((n) => {
              if (n.id === data.node.id) {
                // Preserve original position, only update other fields
                return { ...data.node, position: n.position };
              }
              return n;
            })
          );
        }
        break;

      case 'node_deleted':
        if (data.node) {
          setNodes((prev) => {
            const node = prev.find((n) => n.id === data.node.id);
            if (node) {
              addToast(`Node deleted: ${node.name}`, 'error');
            }
            return prev.filter((n) => n.id !== data.node.id);
          });
        }
        break;
    }
  }, [addToast]);

  // WebSocket callbacks wrapped in useCallback to prevent reconnections
  const handleConnect = useCallback(() => {
    if (!hasConnectedOnce) {
      addToast('Connected to cluster', 'success');
      setHasConnectedOnce(true);
    }
  }, [hasConnectedOnce, addToast]);

  const handleDisconnect = useCallback(() => {
    if (hasConnectedOnce) {
      addToast('Connection lost - attempting to reconnect', 'warning');
    }
  }, [hasConnectedOnce, addToast]);

  const handleError = useCallback(() => {
    // Don't show error toasts for connection issues
  }, []);

  // WebSocket connection
  const { isConnected } = useWebSocket({
    url: 'ws://localhost:8000/ws',
    onMessage: handleWebSocketMessage,
    onConnect: handleConnect,
    onDisconnect: handleDisconnect,
    onError: handleError,
  });

  // Initial data load
  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        setError(null);

        const [nodesData, podsData] = await Promise.all([
          fetchNodes(),
          fetchPods(),
        ]);

        setNodes(nodesData);
        setRawPods(podsData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data');
        console.error('Error loading cluster data:', err);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  // Extract unique namespaces and filter pods
  const namespaces = useMemo(() => {
    const uniqueNamespaces = new Set(pods.map(p => p.namespace));
    return Array.from(uniqueNamespaces).sort();
  }, [pods]);

  const filteredPods = useMemo(() => {
    if (selectedNamespace === 'all') {
      return pods;
    }
    return pods.filter(p => p.namespace === selectedNamespace);
  }, [pods, selectedNamespace]);

  if (loading && nodes.length === 0) {
    return (
      <div className="loading">
        <div className="spinner"></div>
        <p>Loading cluster data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="error">
        <h2>Error</h2>
        <p>{error}</p>
        <p>Make sure the backend is running on port 8000</p>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="header">
        <h1>Observatory</h1>
        <div className="stats">
          <span className="namespace-selector">
            <FontAwesomeIcon icon={faLayerGroup} />
            <select
              value={selectedNamespace}
              onChange={(e) => setSelectedNamespace(e.target.value)}
              className="namespace-dropdown"
            >
              <option value="all">All Namespaces ({namespaces.length})</option>
              {namespaces.map(ns => (
                <option key={ns} value={ns}>
                  {ns}
                </option>
              ))}
            </select>
          </span>
          <span>
            <FontAwesomeIcon icon={faServer} /> {nodes.length} nodes
          </span>
          <span>
            <FontAwesomeIcon icon={faCube} /> {filteredPods.length}/{pods.length} pods
          </span>
          <span className={isConnected ? 'status-connected' : 'status-disconnected'}>
            <FontAwesomeIcon icon={isConnected ? faCircleDot : faCircle} />
            {isConnected ? ' Live' : ' Offline'}
          </span>
        </div>
      </header>
      <div className="scene-container">
        <Scene
          nodes={nodes}
          pods={filteredPods}
          onPodClick={handlePodClick}
          onNodeClick={handleNodeClick}
        />
      </div>
      <ToastContainer toasts={toasts} onRemove={removeToast} />
      <DetailPanel
        resource={selectedResource}
        resourceType={selectedResourceType}
        onClose={handleDetailPanelClose}
      />
    </div>
  );
}
