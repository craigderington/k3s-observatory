export interface Position {
  x: number;
  y: number;
  z: number;
}

export interface ResourceUsage {
  used: number;
  total: number;
}

export interface Node {
  id: string;
  name: string;
  status: string;
  cpu: ResourceUsage;
  memory: ResourceUsage;
  pods: string[];
  labels: Record<string, string>;
  position: Position;
}

export interface Container {
  name: string;
  status: string;
  restarts: number;
  type: string;        // "main", "sidecar", or "init"
  cpu: number;         // millicores
  memory: number;      // MB
}

export interface Pod {
  id: string;
  name: string;
  namespace: string;
  status: string;
  nodeName: string;
  containers: Container[];
  createdAt: string;
  position: Position;
  cpu: number;         // total CPU usage in millicores
  memory: number;      // total memory usage in MB
}

// Metrics update event
export interface MetricsUpdate {
  type: 'metrics_update';
  pods: {
    podId: string;
    name: string;
    namespace: string;
    totalCpu: number;
    totalMemory: number;
    containers: {
      name: string;
      cpu: number;
      memory: number;
    }[];
    timestamp: string;
  }[];
  timestamp: string;
}
