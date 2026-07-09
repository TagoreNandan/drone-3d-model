import { useEffect, useRef, useState, type RefObject } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls } from "@react-three/drei";
import * as THREE from "three";
import Drone from "@/components/drone/Drone";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface TelemetryData {
  roll: number;
  pitch: number;
  yaw: number;
  altitude_m: number;
  speed_mps: number;
  battery_pct: number;
  ts: string;
}

function SceneController({ telemetryRef }: { telemetryRef: RefObject<TelemetryData | null> }) {
  const droneRef = useRef<THREE.Group>(null);

  useFrame(() => {
    if (!droneRef.current || !telemetryRef.current) return;
    const { roll, pitch, yaw } = telemetryRef.current;

    const rollRad = (roll * Math.PI) / 180;
    const pitchRad = (pitch * Math.PI) / 180;
    const yawRad = (yaw * Math.PI) / 180;
    // Axis mapping:
    // Yaw rotates around Y (Up).
    // Pitch rotates around X (Right/Left tilting).
    // Roll rotates around Z (Forward/Backward tilting).
    // Order YXZ corresponds to standard aerospace rotation sequences in Three.js coordinates.
    const euler = new THREE.Euler(pitchRad, yawRad, rollRad, "YXZ");
    droneRef.current.quaternion.setFromEuler(euler);
  });

  return <Drone ref={droneRef} />;
}

export default function DroneView() {
  const telemetryRef = useRef<TelemetryData | null>(null);
  const [metrics, setMetrics] = useState<TelemetryData | null>(null);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws/mock/orientation");
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as TelemetryData;
        telemetryRef.current = data;
      } catch (err) {
        console.error("Failed to parse telemetry message:", err);
      }
    };

    // Modest rate sync for metrics state to prevent React render storms at 50Hz.
    const interval = setInterval(() => {
      if (telemetryRef.current) {
        setMetrics(telemetryRef.current);
      }
    }, 200);

    return () => {
      clearInterval(interval);
      ws.onmessage = null;
      ws.onclose = null;
      ws.onerror = null;
      ws.onopen = null;
      ws.close();
    };
  }, []);

  return (
    <div className="w-full h-full relative">
      {/* Overlay displaying telemetry metrics */}
      {metrics && (
        <div className="absolute top-4 left-4 z-10 w-52 pointer-events-none">
          <Card className="pointer-events-auto shadow-md">
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-semibold">Telemetry</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-2">
              <div className="flex justify-between items-center text-xs">
                <span className="text-muted-foreground">Altitude:</span>
                <span className="font-mono font-semibold">{metrics.altitude_m.toFixed(1)} m</span>
              </div>
              <div className="flex justify-between items-center text-xs">
                <span className="text-muted-foreground">Speed:</span>
                <span className="font-mono font-semibold">{metrics.speed_mps.toFixed(1)} m/s</span>
              </div>
              <div className="flex justify-between items-center text-xs">
                <span className="text-muted-foreground">Battery:</span>
                <span className="font-mono font-semibold">{Math.round(metrics.battery_pct)} %</span>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <Canvas camera={{ position: [2, 2, 2], fov: 50 }}>
        <ambientLight intensity={0.5} />
        <directionalLight position={[10, 10, 10]} intensity={1.0} />
        <SceneController telemetryRef={telemetryRef} />
        <OrbitControls />
      </Canvas>
    </div>
  );
}
