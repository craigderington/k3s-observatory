import { Canvas } from '@react-three/fiber';
import { OrbitControls, Grid } from '@react-three/drei';
import { Node, Pod } from '../types';
import NodeSphere from './NodeSphere';
import PodSphere from './PodSphere';

interface SceneProps {
  nodes: Node[];
  pods: Pod[];
  onPodClick?: (pod: Pod) => void;
  onNodeClick?: (node: Node) => void;
}

export default function Scene({ nodes, pods, onPodClick, onNodeClick }: SceneProps) {
  return (
    <Canvas camera={{ position: [15, 15, 15], fov: 60 }}>
      <color attach="background" args={['#0a0a0a']} />

      {/* Lighting */}
      <ambientLight intensity={0.5} />
      <pointLight position={[10, 10, 10]} intensity={1} />
      <pointLight position={[-10, -10, -10]} intensity={0.5} />

      {/* Grid for reference */}
      <Grid args={[50, 50]} cellSize={1} cellColor="#333333" sectionColor="#666666" />

      {/* Render nodes */}
      {nodes.map((node) => (
        <NodeSphere key={node.id} node={node} onClick={onNodeClick} />
      ))}

      {/* Render pods */}
      {pods.map((pod) => (
        <PodSphere key={pod.id} pod={pod} onClick={onPodClick} />
      ))}

      {/* Camera controls */}
      <OrbitControls />
    </Canvas>
  );
}
