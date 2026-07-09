import { forwardRef } from "react";
import type { ComponentProps } from "react";
import type { Group } from "three";

const Drone = forwardRef<Group, ComponentProps<"group">>((props, ref) => {
  return (
    <group ref={ref} {...props}>
      {/* Central body */}
      <mesh position={[0, 0, 0]}>
        <boxGeometry args={[0.3, 0.1, 0.3]} />
        <meshStandardMaterial color="#1e3a8a" roughness={0.5} />
      </mesh>

      {/* Diagonal arm crossbars */}
      <mesh rotation={[0, Math.PI / 4, 0]}>
        <boxGeometry args={[1.2, 0.04, 0.04]} />
        <meshStandardMaterial color="#333333" />
      </mesh>
      <mesh rotation={[0, -Math.PI / 4, 0]}>
        <boxGeometry args={[1.2, 0.04, 0.04]} />
        <meshStandardMaterial color="#333333" />
      </mesh>

      {/* Motors and Propellers at the 4 ends (radius from center = 0.6) */}
      {/* Front Left: x = -0.424, z = -0.424 */}
      <group position={[-0.424, 0.04, -0.424]}>
        <mesh>
          <cylinderGeometry args={[0.03, 0.03, 0.08, 16]} />
          <meshStandardMaterial color="#888888" metalness={0.8} roughness={0.2} />
        </mesh>
        <mesh position={[0, 0.05, 0]}>
          <cylinderGeometry args={[0.15, 0.15, 0.005, 16]} />
          <meshStandardMaterial color="#f59e0b" opacity={0.7} transparent />
        </mesh>
      </group>

      {/* Front Right: x = 0.424, z = -0.424 */}
      <group position={[0.424, 0.04, -0.424]}>
        <mesh>
          <cylinderGeometry args={[0.03, 0.03, 0.08, 16]} />
          <meshStandardMaterial color="#888888" metalness={0.8} roughness={0.2} />
        </mesh>
        <mesh position={[0, 0.05, 0]}>
          <cylinderGeometry args={[0.15, 0.15, 0.005, 16]} />
          <meshStandardMaterial color="#f59e0b" opacity={0.7} transparent />
        </mesh>
      </group>

      {/* Back Left: x = -0.424, z = 0.424 */}
      <group position={[-0.424, 0.04, 0.424]}>
        <mesh>
          <cylinderGeometry args={[0.03, 0.03, 0.08, 16]} />
          <meshStandardMaterial color="#888888" metalness={0.8} roughness={0.2} />
        </mesh>
        <mesh position={[0, 0.05, 0]}>
          <cylinderGeometry args={[0.15, 0.15, 0.005, 16]} />
          <meshStandardMaterial color="#f59e0b" opacity={0.7} transparent />
        </mesh>
      </group>

      {/* Development orientation marker. Indicates the forward (-Z) direction of the placeholder drone. Remove once the production GLB model is used. */}
      <mesh position={[0, 0.05, -0.34]} rotation={[Math.PI / 2, 0, 0]}>
        <coneGeometry args={[0.05, 0.14, 16]} />
        <meshStandardMaterial color="#ef4444" />
      </mesh>

      {/* Back Right: x = 0.424, z = 0.424 */}
      <group position={[0.424, 0.04, 0.424]}>
        <mesh>
          <cylinderGeometry args={[0.03, 0.03, 0.08, 16]} />
          <meshStandardMaterial color="#888888" metalness={0.8} roughness={0.2} />
        </mesh>
        <mesh position={[0, 0.05, 0]}>
          <cylinderGeometry args={[0.15, 0.15, 0.005, 16]} />
          <meshStandardMaterial color="#f59e0b" opacity={0.7} transparent />
        </mesh>
      </group>
    </group>
  );
});

Drone.displayName = "Drone";

export default Drone;
