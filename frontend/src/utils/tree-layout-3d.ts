import type { FamilyTreeData, Person, Relationship } from "@/types/person";

export interface Node3D {
  person: Person;
  x: number;
  y: number;
  z: number;
  generation: number;
}

export interface Layout3D {
  nodes: Node3D[];
  edges: Relationship[];
}

const X_GAP = 2.4;
const Y_GAP = 2.8;
const Z_JITTER = 0.35;

/** Hierarchical 3D layout: generations on Y, siblings on X, slight Z depth. */
export function layoutTree3D(data: FamilyTreeData): Layout3D {
  if (data.persons.length === 0) {
    return { nodes: [], edges: [] };
  }

  const generations = new Map<string, number>();
  for (const person of data.persons) {
    generations.set(person.id, 0);
  }

  const parentRels = data.relationships.filter((r) => r.type === "parent");
  const spouseRels = data.relationships.filter((r) => r.type === "spouse");

  let changed = true;
  while (changed) {
    changed = false;
    for (const rel of parentRels) {
      const parentGen = generations.get(rel.from_person_id) ?? 0;
      const childGen = generations.get(rel.to_person_id) ?? 0;
      const next = parentGen + 1;
      if (next > childGen) {
        generations.set(rel.to_person_id, next);
        changed = true;
      }
    }
  }

  for (const rel of spouseRels) {
    const a = generations.get(rel.from_person_id) ?? 0;
    const b = generations.get(rel.to_person_id) ?? 0;
    const aligned = Math.min(a, b);
    generations.set(rel.from_person_id, aligned);
    generations.set(rel.to_person_id, aligned);
  }

  const byGeneration = new Map<number, Person[]>();
  for (const person of data.persons) {
    const gen = generations.get(person.id) ?? 0;
    const bucket = byGeneration.get(gen) ?? [];
    bucket.push(person);
    byGeneration.set(gen, bucket);
  }

  const nodes: Node3D[] = [];
  const sortedGens = [...byGeneration.keys()].sort((a, b) => a - b);

  for (const gen of sortedGens) {
    const row = byGeneration.get(gen) ?? [];
    row.sort((a, b) => displayName(a).localeCompare(displayName(b)));
    const offset = (row.length - 1) / 2;
    row.forEach((person, index) => {
      const hash = hashId(person.id);
      nodes.push({
        person,
        generation: gen,
        x: (index - offset) * X_GAP,
        y: -gen * Y_GAP,
        z: ((hash % 100) / 100 - 0.5) * 2 * Z_JITTER,
      });
    });
  }

  const ids = new Set(data.persons.map((p) => p.id));
  return {
    nodes,
    edges: data.relationships.filter(
      (rel) => ids.has(rel.from_person_id) && ids.has(rel.to_person_id),
    ),
  };
}

export function displayName(person: Person): string {
  return person.last_name ? `${person.first_name} ${person.last_name}` : person.first_name;
}

function hashId(id: string): number {
  let h = 0;
  for (let i = 0; i < id.length; i++) {
    h = (h * 31 + id.charCodeAt(i)) | 0;
  }
  return Math.abs(h);
}

export function findNode3D(nodes: Node3D[], id: string): Node3D | undefined {
  return nodes.find((n) => n.person.id === id);
}
