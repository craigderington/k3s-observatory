import { useState, useEffect } from 'react';
import Scene from './components/Scene';
import { fetchNodes, fetchPods } from './services/api';
import { Node, Pod } from './types';
import './App.css';

export default function App() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [pods, setPods] = useState<Pod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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
        setPods(podsData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data');
        console.error('Error loading cluster data:', err);
      } finally {
        setLoading(false);
      }
    };

    loadData();

    // Refresh data every 5 seconds
    const interval = setInterval(loadData, 5000);
    return () => clearInterval(interval);
  }, []);

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
          <span>{nodes.length} nodes</span>
          <span>{pods.length} pods</span>
        </div>
      </header>
      <div className="scene-container">
        <Scene nodes={nodes} pods={pods} />
      </div>
    </div>
  );
}
