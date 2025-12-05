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

export async function describePod(namespace: string, name: string): Promise<string> {
  const response = await fetch(`${API_BASE}/pods/describe?namespace=${encodeURIComponent(namespace)}&name=${encodeURIComponent(name)}`);
  if (!response.ok) {
    throw new Error('Failed to describe pod');
  }
  return response.text();
}

export async function describeNode(name: string): Promise<string> {
  const response = await fetch(`${API_BASE}/nodes/describe?name=${encodeURIComponent(name)}`);
  if (!response.ok) {
    throw new Error('Failed to describe node');
  }
  return response.text();
}

export async function getPodLogs(namespace: string, name: string, container?: string): Promise<string> {
  let url = `${API_BASE}/pods/logs?namespace=${encodeURIComponent(namespace)}&name=${encodeURIComponent(name)}`;
  if (container) {
    url += `&container=${encodeURIComponent(container)}`;
  }
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error('Failed to get pod logs');
  }
  return response.text();
}

export interface Metrics {
  name: string;
  namespace?: string;
  cpuUsage: number;    // in millicores
  memoryUsage: number; // in MB
  timestamp: string;
}

export async function getPodMetrics(namespace: string, name: string): Promise<Metrics> {
  const response = await fetch(`${API_BASE}/pods/metrics?namespace=${encodeURIComponent(namespace)}&name=${encodeURIComponent(name)}`);
  if (!response.ok) {
    throw new Error('Failed to get pod metrics');
  }
  return response.json();
}

export async function getNodeMetrics(name: string): Promise<Metrics> {
  const response = await fetch(`${API_BASE}/nodes/metrics?name=${encodeURIComponent(name)}`);
  if (!response.ok) {
    throw new Error('Failed to get node metrics');
  }
  return response.json();
}
