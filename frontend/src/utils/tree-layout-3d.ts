import type { FamilyTreeData, Person, Relationship } from "@/types/person";
import { peopleByGeneration } from "@/utils/tree-layout";

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

  const rows = peopleByGeneration(data);
  const nodes: Node3D[] = [];

  rows.forEach((row, gen) => {
    const offset = (row.length - 1) / 2;
    row.forEach((person, index) => {
      nodes.push({
        person,
        generation: gen,
        x: (index - offset) * X_GAP,
        y: -gen * Y_GAP,
        z: 0,
      });
    });
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
