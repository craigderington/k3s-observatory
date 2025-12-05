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
}
