import type { FamilyTreeData, Person, Relationship } from "@/types/person";
import { arrangeTree } from "@/utils/tree-layout";

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

/** Hierarchical 3D layout: generations on Y, spouses beside each other, children below. */
export function layoutTree3D(data: FamilyTreeData): Layout3D {
  if (data.persons.length === 0) {
    return { nodes: [], edges: [] };
  }

  const { rows, slots } = arrangeTree(data);
  const slotValues = [...slots.values()];
  const mid = (Math.min(...slotValues) + Math.max(...slotValues)) / 2;
  const nodes: Node3D[] = [];

  rows.forEach((row, gen) => {
    for (const person of row) {
      nodes.push({
        person,
        generation: gen,
        x: ((slots.get(person.id) ?? mid) - mid) * X_GAP,
        y: -gen * Y_GAP,
        z: 0,
      });
    }
  });

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
