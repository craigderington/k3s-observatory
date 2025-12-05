import { Node, Pod } from '../types';

const API_BASE = '/api';

export async function fetchNodes(): Promise<Node[]> {
  const response = await fetch(`${API_BASE}/nodes`);
  if (!response.ok) {
    throw new Error('Failed to fetch nodes');
  }
  return response.json();
}

export async function fetchPods(): Promise<Pod[]> {
  const response = await fetch(`${API_BASE}/pods`);
  if (!response.ok) {
    throw new Error('Failed to fetch pods');
  }
  return response.json();
}

export async function checkHealth(): Promise<{ status: string; service: string }> {
  const response = await fetch(`${API_BASE}/health`);
  if (!response.ok) {
    throw new Error('Health check failed');
  }
  return response.json();
}
