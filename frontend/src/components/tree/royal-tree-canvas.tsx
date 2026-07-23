"use client";

import { Billboard, Float, GradientTexture, Html, Line, OrbitControls } from "@react-three/drei";
import { Canvas } from "@react-three/fiber";
import { Suspense, useEffect, useMemo, useState } from "react";
import * as THREE from "three";

import type { FamilyTreeData, Person } from "@/types/person";
import { displayName, findNode3D, layoutTree3D, type Node3D } from "@/utils/tree-layout-3d";

interface RoyalTreeCanvasProps {
  data: FamilyTreeData;
  selectedId?: string | null;
  onSelect: (person: Person) => void;
}

const GOLD = "#d4af37";
const GOLD_BRIGHT = "#f0c75e";
const SEA_DEEP = "#0a3d42";
const SEA_FACE = "#124a50";
const INK = "#04151a";

function PersonNode({
  node,
  selected,
  onSelect,
}: {
  node: Node3D;
  selected: boolean;
  onSelect: () => void;
}) {
  const name = displayName(node.person);
  const initials = initialsFrom(name);
  const photoUrl = node.person.photo_url?.trim() || null;
  const rimColor = selected ? GOLD_BRIGHT : GOLD;

  return (
    <Float speed={0.6} rotationIntensity={0.05} floatIntensity={0.12}>
      <group position={[node.x, node.y, node.z]}>
        <Billboard follow>
          <group
            onClick={(e) => {
              e.stopPropagation();
              onSelect();
            }}
            onPointerOver={(e) => {
              e.stopPropagation();
              document.body.style.cursor = "pointer";
            }}
            onPointerOut={() => {
              document.body.style.cursor = "auto";
            }}
          >
            {/* Soft outer glow when selected */}
            {selected && (
              <mesh position={[0, 0, -0.03]}>
                <circleGeometry args={[0.42, 48]} />
                <meshBasicMaterial color={GOLD_BRIGHT} transparent opacity={0.22} />
              </mesh>
            )}

            {/* Gold rim (raised seal edge) */}
            <mesh position={[0, 0, 0.01]}>
              <torusGeometry args={[0.3, 0.035, 16, 48]} />
              <meshStandardMaterial
                color={rimColor}
                emissive={rimColor}
                emissiveIntensity={selected ? 0.35 : 0.12}
                metalness={0.85}
                roughness={0.28}
              />
            </mesh>

            {/* Inner face: photo or monogram */}
            <MedallionFace photoUrl={photoUrl} initials={initials} selected={selected} />

            {/* Thin back plate for depth */}
            <mesh position={[0, 0, -0.02]}>
              <circleGeometry args={[0.31, 48]} />
              <meshStandardMaterial color={INK} metalness={0.2} roughness={0.7} />
            </mesh>
          </group>
        </Billboard>

        <Html distanceFactor={9} position={[0, -0.58, 0]} center>
          <NamePlaque name={name} selected={selected} onSelect={onSelect} />
        </Html>
      </group>
    </Float>
  );
}

/** Fixed-size name plaque: 2-line clamp + smaller type for long names. */
function NamePlaque({
  name,
  selected,
  onSelect,
}: {
  name: string;
  selected: boolean;
  onSelect: () => void;
}) {
  const long = name.length > 26;
  const veryLong = name.length > 38;

  return (
    <button
      type="button"
      title={name}
      aria-label={name}
      onClick={(e) => {
        e.stopPropagation();
        onSelect();
      }}
      className={
        selected
          ? "box-border flex h-[2.6rem] w-[9.75rem] cursor-pointer items-center justify-center rounded-md bg-[#04151a]/92 px-2 py-1 text-center font-medium text-amber-100 shadow-lg ring-1 ring-amber-400/50 backdrop-blur-sm"
          : "box-border flex h-[2.6rem] w-[9.75rem] cursor-pointer items-center justify-center rounded-md bg-[#04151a]/85 px-2 py-1 text-center font-medium text-teal-50 shadow-lg ring-1 ring-teal-400/25 backdrop-blur-sm"
      }
    >
      <span
        className={
          veryLong
            ? "line-clamp-2 w-full break-words text-[9px] leading-tight"
            : long
              ? "line-clamp-2 w-full break-words text-[10px] leading-snug"
              : "line-clamp-2 w-full break-words text-[11px] leading-snug"
        }
      >
        {name}
      </span>
    </button>
  );
}

function MedallionFace({
  photoUrl,
  initials,
  selected,
}: {
  photoUrl: string | null;
  initials: string;
  selected: boolean;
}) {
  const texture = useCoverTexture(photoUrl);

  if (texture) {
    return (
      <mesh>
        <circleGeometry args={[0.275, 48]} />
        <meshStandardMaterial
          map={texture}
          roughness={0.55}
          metalness={0.05}
          emissive={selected ? "#3a2a10" : "#000000"}
          emissiveIntensity={selected ? 0.15 : 0}
        />
      </mesh>
    );
  }

  return (
    <group>
      <mesh>
        <circleGeometry args={[0.275, 48]} />
        <meshStandardMaterial
          color={selected ? "#1a5c58" : SEA_FACE}
          emissive={SEA_DEEP}
          emissiveIntensity={0.25}
          metalness={0.35}
          roughness={0.45}
        />
      </mesh>
      <Html transform distanceFactor={6} position={[0, 0, 0.02]} center>
        <div
          className="pointer-events-none select-none font-display text-[22px] font-semibold tracking-wide text-amber-100/95"
          style={{ textShadow: "0 1px 2px rgba(0,0,0,0.55)" }}
        >
          {initials}
        </div>
      </Html>
    </group>
  );
}

/** Loads a portrait texture and crops it like CSS object-fit: cover. */
function useCoverTexture(url: string | null): THREE.Texture | null {
  const [texture, setTexture] = useState<THREE.Texture | null>(null);

  useEffect(() => {
    if (!url) {
      setTexture(null);
      return;
    }

    let cancelled = false;
    let loaded: THREE.Texture | null = null;
    const loader = new THREE.TextureLoader();
    loader.setCrossOrigin("anonymous");

    loader.load(
      url,
      (tex) => {
        if (cancelled) {
          tex.dispose();
          return;
        }
        tex.colorSpace = THREE.SRGBColorSpace;
        tex.wrapS = THREE.ClampToEdgeWrapping;
        tex.wrapT = THREE.ClampToEdgeWrapping;

        const image = tex.image as { width?: number; height?: number } | undefined;
        const width = image?.width ?? 1;
        const height = image?.height ?? 1;
        const aspect = width / height;

        // object-fit: cover onto the UV square (circle clips the rest).
        if (aspect > 1) {
          tex.repeat.set(1 / aspect, 1);
          tex.offset.set((1 - 1 / aspect) / 2, 0);
        } else {
          tex.repeat.set(1, aspect);
          tex.offset.set(0, (1 - aspect) / 2);
        }

        tex.needsUpdate = true;
        loaded = tex;
        setTexture(tex);
      },
      undefined,
      () => {
        if (!cancelled) setTexture(null);
      },
    );

    return () => {
      cancelled = true;
      loaded?.dispose();
      setTexture(null);
    };
  }, [url]);

  return texture;
}

function initialsFrom(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return `${parts[0][0] ?? ""}${parts[parts.length - 1][0] ?? ""}`.toUpperCase();
}

function EdgeLines({ nodes, edges }: { nodes: Node3D[]; edges: FamilyTreeData["relationships"] }) {
  return (
    <>
      {edges.map((edge) => {
        const from = findNode3D(nodes, edge.from_person_id);
        const to = findNode3D(nodes, edge.to_person_id);
        if (!from || !to) return null;
        const points = [
          new THREE.Vector3(from.x, from.y, from.z),
          new THREE.Vector3(to.x, to.y, to.z),
        ];
        const isSpouse = edge.type === "spouse";
        return (
          <Line
            key={edge.id}
            points={points}
            color={isSpouse ? GOLD : "#2a8f86"}
            lineWidth={isSpouse ? 1.1 : 1.6}
            dashed={isSpouse}
            dashSize={0.12}
            gapSize={0.1}
            transparent
            opacity={isSpouse ? 0.7 : 0.8}
          />
        );
      })}
    </>
  );
}

/** Soft Aceh sea depth: dusk gold above, deep teal water below. */
function SeaDepthSky() {
  return (
    <group>
      <mesh scale={[-1, 1, 1]}>
        <sphereGeometry args={[70, 48, 32]} />
        <meshBasicMaterial side={THREE.BackSide} depthWrite={false}>
          <GradientTexture
            attach="map"
            stops={[0, 0.35, 0.62, 1]}
            colors={["#1a3538", "#0c3d44", "#062a30", "#04151a"]}
            size={256}
          />
        </meshBasicMaterial>
      </mesh>

      {/* Faint gold dusk wash near the zenith */}
      <mesh position={[0, 28, -8]} rotation={[0.35, 0, 0]}>
        <circleGeometry args={[22, 48]} />
        <meshBasicMaterial
          color={GOLD}
          transparent
          opacity={0.07}
          depthWrite={false}
          side={THREE.DoubleSide}
        />
      </mesh>

      {/* Soft teal bloom low in the field — “sea glow” under the tree */}
      <mesh position={[0, -10, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <circleGeometry args={[18, 48]} />
        <meshBasicMaterial
          color="#1a9b8e"
          transparent
          opacity={0.06}
          depthWrite={false}
          side={THREE.DoubleSide}
        />
      </mesh>
    </group>
  );
}

function Scene({ data, selectedId, onSelect }: RoyalTreeCanvasProps) {
  const layout = useMemo(() => layoutTree3D(data), [data]);

  return (
    <>
      <color attach="background" args={["#062a30"]} />
      <fog attach="fog" args={["#062a30", 18, 46]} />
      <SeaDepthSky />

      <ambientLight intensity={0.58} />
      <directionalLight position={[8, 12, 6]} intensity={1.1} color="#f5e6c8" />
      <pointLight position={[-6, 4, -4]} intensity={0.4} color="#1a9b8e" />
      <pointLight position={[2, 10, 4]} intensity={0.28} color={GOLD} />

      <EdgeLines nodes={layout.nodes} edges={layout.edges} />
      {layout.nodes.map((node) => (
        <PersonNode
          key={node.person.id}
          node={node}
          selected={selectedId === node.person.id}
          onSelect={() => onSelect(node.person)}
        />
      ))}

      <OrbitControls
        makeDefault
        enableDamping
        dampingFactor={0.08}
        minDistance={4}
        maxDistance={36}
        maxPolarAngle={Math.PI * 0.9}
      />
    </>
  );
}

/** Full-viewport WebGL royal tree. Dynamically import this component. */
export function RoyalTreeCanvas(props: RoyalTreeCanvasProps) {
  if (typeof window !== "undefined") {
    try {
      const canvas = document.createElement("canvas");
      const gl = canvas.getContext("webgl") || canvas.getContext("experimental-webgl");
      if (!gl) {
        return (
          <div className="aceh-sea-bg grid h-full place-items-center px-6 text-center text-teal-100/80">
            WebGL is not available in this browser. Try a recent Chrome, Firefox, or Edge.
          </div>
        );
      }
    } catch {
      // continue — Canvas will surface errors
    }
  }

  return (
    <Canvas
      className="h-full w-full touch-none"
      camera={{ position: [0, 2, 12], fov: 50 }}
      dpr={[1, 1.75]}
    >
      <Suspense fallback={null}>
        <Scene {...props} />
      </Suspense>
    </Canvas>
  );
}
