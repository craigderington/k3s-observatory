import { useRef, useState } from 'react';
import { Mesh } from 'three';
import { Node } from '../types';
import { Text } from '@react-three/drei';

interface NodeSphereProps {
  node: Node;
}

export default function NodeSphere({ node }: NodeSphereProps) {
  const meshRef = useRef<Mesh>(null);
  const [hovered, setHovered] = useState(false);

  // Color based on status
  const color = node.status === 'Ready' ? '#00ff00' : '#ff0000';

  return (
    <group position={[node.position.x, node.position.y, node.position.z]}>
      {/* Node sphere */}
      <mesh
        ref={meshRef}
        onPointerOver={() => setHovered(true)}
        onPointerOut={() => setHovered(false)}
        scale={hovered ? 1.2 : 1}
      >
        <sphereGeometry args={[1, 32, 32]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={hovered ? 0.5 : 0.2}
        />
      </mesh>

      {/* Node label */}
      <Text
        position={[0, 1.5, 0]}
        fontSize={0.3}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {node.name}
      </Text>

      {/* Status indicator */}
      {hovered && (
        <Text
          position={[0, -1.5, 0]}
          fontSize={0.2}
          color="#aaaaaa"
          anchorX="center"
          anchorY="middle"
        >
          {node.status}
        </Text>
      )}
    </group>
  );
}
