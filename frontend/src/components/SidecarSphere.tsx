import { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import { Mesh } from 'three';
import { Container } from '../types';

interface SidecarSphereProps {
  container: Container;
  podScale: number;
  index: number;
  total: number;
}

export default function SidecarSphere({ container, podScale, index, total }: SidecarSphereProps) {
  const meshRef = useRef<Mesh>(null);

  // Calculate orbit angle for this sidecar
  const angle = (index * 2 * Math.PI) / total;

  // Orbit radius proportional to pod size
  const orbitRadius = podScale * 1.5;

  // Sidecar size is 20-25% of main pod
  const sidecarScale = podScale * 0.25;

  // Rotate sidecars around pod (animation)
  useFrame(({ clock }) => {
    if (meshRef.current) {
      const time = clock.getElapsedTime();
      const orbitSpeed = 0.5; // Slow rotation
      const currentAngle = angle + time * orbitSpeed;

      meshRef.current.position.x = orbitRadius * Math.cos(currentAngle);
      meshRef.current.position.z = orbitRadius * Math.sin(currentAngle);
    }
  });

  // Color based on container status
  const getColor = () => {
    switch (container.status) {
      case 'Running':
        return '#00ff00';
      case 'Waiting':
        return '#0088ff';
      case 'Terminated':
        return '#ff0000';
      default:
        return '#ffff00';
    }
  };

  // CPU heat glow for sidecars
  const getCPUHeatColor = () => {
    if (!container.cpu || container.cpu === 0) {
      return '#000000';
    }

    // Sidecars typically use less CPU, so use 500m as max
    const cpuPercent = Math.min(container.cpu / 500, 1.0);

    if (cpuPercent < 0.5) {
      return `rgb(${Math.floor(cpuPercent * 2 * 255)}, ${Math.floor(cpuPercent * 2 * 200)}, 255)`;
    } else {
      return `rgb(255, ${Math.floor(255 - (cpuPercent - 0.5) * 2 * 90)}, 0)`;
    }
  };

  const color = getColor();

  return (
    <mesh
      ref={meshRef}
      position={[orbitRadius * Math.cos(angle), 0, orbitRadius * Math.sin(angle)]}
      scale={sidecarScale}
    >
      <sphereGeometry args={[1, 12, 12]} />
      <meshStandardMaterial
        color={color}
        emissive={getCPUHeatColor()}
        emissiveIntensity={container.cpu > 0 ? 0.3 : 0.1}
        opacity={0.8}
        transparent={true}
      />
    </mesh>
  );
}
