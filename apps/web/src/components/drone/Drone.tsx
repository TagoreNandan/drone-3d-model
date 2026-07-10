import { useGLTF } from "@react-three/drei";
import { useEffect } from "react";
import * as THREE from "three";

export default function Drone() {
  const { scene } = useGLTF("/models/drone.glb");
  scene.rotation.y = Math.PI;

  useEffect(() => {
    const box = new THREE.Box3().setFromObject(scene);
    const size = box.getSize(new THREE.Vector3());
    const center = box.getCenter(new THREE.Vector3());

    console.log("Model size:", size);
    console.log("Model center:", center);

    scene.position.sub(center);
  }, [scene]);

  return <primitive object={scene} scale={10} />;
}

useGLTF.preload("/models/drone.glb");