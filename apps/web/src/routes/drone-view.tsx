import { useEffect, useRef, useState, type RefObject } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls, useGLTF } from "@react-three/drei";
import * as THREE from "three";
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

function SceneController({
  telemetryRef,
}: {
  telemetryRef: React.RefObject<TelemetryData | null>;
}) {
  const { scene } = useGLTF("/models/drone.glb");

  useEffect(() => {
    scene.rotation.y = Math.PI;

    const box = new THREE.Box3().setFromObject(scene);
    const center = box.getCenter(new THREE.Vector3());
    scene.position.sub(center);
  }, [scene]);

  useFrame(() => {
    if (!telemetryRef.current) return;

    const { roll, pitch, yaw } = telemetryRef.current;

    scene.rotation.set(
      THREE.MathUtils.degToRad(pitch),
      THREE.MathUtils.degToRad(yaw),
      THREE.MathUtils.degToRad(roll)
    );
  });

  return <primitive object={scene} scale={10} />;
}

export default function DroneView() {
  const telemetryRef = useRef<TelemetryData | null>(null);
  const [metrics, setMetrics] = useState<TelemetryData | null>(null);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws/mock/orientation");
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as TelemetryData;

        console.log(data);

        telemetryRef.current = data;
      } catch (err) {
        console.error(err);
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

      <Canvas camera={{ position: [5, 4, 8], fov: 45 }}>
        <ambientLight intensity={5} />
        <directionalLight
          position={[5, 5, 5]}
          intensity={8}
        />
        <directionalLight
          position={[-5, 5, -5]}
          intensity={6}
        />
        <hemisphereLight
          intensity={3}
          groundColor="gray"
        />
        <SceneController telemetryRef={telemetryRef} />
        <OrbitControls
          target={[0, 0.8, 0]}
          enablePan={false}
          enableDamping
          dampingFactor={0.08}
        />
      </Canvas>
    </div>
  );
}
