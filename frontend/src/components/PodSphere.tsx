import { useRef, useState } from 'react';
import { Mesh } from 'three';
import { Pod } from '../types';
import { Text } from '@react-three/drei';

interface PodSphereProps {
  pod: Pod;
  onClick?: (pod: Pod) => void;
}

export default function PodSphere({ pod, onClick }: PodSphereProps) {
  const meshRef = useRef<Mesh>(null);
  const [hovered, setHovered] = useState(false);

  // Color based on pod status
  const getColor = () => {
    switch (pod.status) {
      case 'Running':
        return '#00ff00';
      case 'Pending':
        return '#0088ff';
      case 'Failed':
        return '#ff0000';
      case 'Succeeded':
        return '#888888';
      default:
        return '#ffff00';
    }
  };

  const color = getColor();

  return (
    <group position={[pod.position.x, pod.position.y, pod.position.z]}>
      {/* Pod sphere */}
      <mesh
        ref={meshRef}
        onPointerOver={() => setHovered(true)}
        onPointerOut={() => setHovered(false)}
        onClick={(e) => {
          e.stopPropagation();
          onClick?.(pod);
        }}
        scale={hovered ? 0.6 : 0.5}
      >
        <sphereGeometry args={[1, 16, 16]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={hovered ? 0.4 : 0.1}
        />
      </mesh>

      {/* Pod label - always visible */}
      <Text
        position={[0, hovered ? 0.8 : 0.7, 0]}
        fontSize={hovered ? 0.15 : 0.1}
        color="white"
        anchorX="center"
        anchorY="middle"
        outlineWidth={0.01}
        outlineColor="#000000"
      >
        {hovered
          ? (pod.name.length > 20 ? pod.name.substring(0, 20) + '...' : pod.name)
          : (pod.name.length > 15 ? pod.name.substring(0, 12) + '...' : pod.name)}
      </Text>

      {/* Namespace (only show on hover) */}
      {hovered && (
        <Text
          position={[0, -0.8, 0]}
          fontSize={0.12}
          color="#aaaaaa"
          anchorX="center"
          anchorY="middle"
          outlineWidth={0.01}
          outlineColor="#000000"
        >
          {pod.namespace}
        </Text>
      )}
    </group>
  );
}
