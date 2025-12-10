import { useRef, useState } from 'react';
import { Mesh } from 'three';
import { Pod } from '../types';
import { Text } from '@react-three/drei';
import SidecarSphere from './SidecarSphere';

interface PodSphereProps {
  pod: Pod;
  onClick?: (pod: Pod) => void;
}

export default function PodSphere({ pod, onClick }: PodSphereProps) {
  const meshRef = useRef<Mesh>(null);
  const [hovered, setHovered] = useState(false);

  // Separate sidecar containers from main container
  const sidecarContainers = pod.containers.filter(c => c.type === 'sidecar');

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

  // Scale based on memory usage
  const getScale = () => {
    const baseScale = 0.5;

    // If no memory data, use base scale
    if (!pod.memory || pod.memory === 0) {
      return baseScale;
    }

    // Logarithmic scaling: 0-1024MB → 0.3-1.0 scale
    const memoryMB = pod.memory;
    const minScale = 0.3;
    const maxScale = 1.0;
    const memoryRange = 1024; // MB

    // Logarithmic scaling for better visual distribution
    const scaleFactor = Math.log10(memoryMB + 1) / Math.log10(memoryRange + 1);
    const scale = minScale + (maxScale - minScale) * scaleFactor;

    return hovered ? scale * 1.2 : scale;
  };

  // CPU heat color gradient
  const getCPUHeatColor = () => {
    // If no CPU data, return no heat (black)
    if (!pod.cpu || pod.cpu === 0) {
      return '#000000';
    }

    // Normalize to 0-1 (1000m = 1 core)
    const cpuPercent = Math.min(pod.cpu / 1000, 1.0);

    if (cpuPercent < 0.25) {
      // Blue to Cyan (0-25%)
      return `rgb(${Math.floor(cpuPercent * 4 * 100)}, ${Math.floor(cpuPercent * 4 * 200)}, 255)`;
    } else if (cpuPercent < 0.5) {
      // Cyan to Yellow (25-50%)
      const t = (cpuPercent - 0.25) * 4;
      return `rgb(${Math.floor(200 + t * 55)}, ${Math.floor(200 + t * 55)}, ${Math.floor(255 - t * 255)})`;
    } else if (cpuPercent < 0.75) {
      // Yellow to Orange (50-75%)
      const t = (cpuPercent - 0.5) * 4;
      return `rgb(255, ${Math.floor(255 - t * 90)}, 0)`;
    } else {
      // Orange to Red (75-100%)
      const t = (cpuPercent - 0.75) * 4;
      return `rgb(255, ${Math.floor(165 - t * 165)}, 0)`;
    }
  };

  // CPU emissive intensity
  const getCPUEmissiveIntensity = () => {
    if (!pod.cpu || pod.cpu === 0) {
      return 0.1;
    }

    // Higher CPU = more intense glow
    const cpuPercent = Math.min(pod.cpu / 1000, 1.0);
    return 0.2 + cpuPercent * 0.6; // Range: 0.2 to 0.8
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
        scale={getScale()}
      >
        <sphereGeometry args={[1, 16, 16]} />
        <meshStandardMaterial
          color={color}
          emissive={getCPUHeatColor()}
          emissiveIntensity={getCPUEmissiveIntensity()}
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

      {/* Sidecar satellites */}
      {sidecarContainers.map((container, index) => (
        <SidecarSphere
          key={container.name}
          container={container}
          podScale={getScale()}
          index={index}
          total={sidecarContainers.length}
        />
      ))}
    </group>
  );
}
