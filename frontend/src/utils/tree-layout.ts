import type { FamilyTreeData, Person, Relationship, TreeLayout, TreeNodePosition } from "@/types/person";

const NODE_WIDTH = 152;
const NODE_HEIGHT = 76;
const H_GAP = 48;
const V_GAP = 96;
const PADDING = 40;

/** Assigns generation rows and x positions for a simple top-down tree layout. */
export function layoutFamilyTree(data: FamilyTreeData): TreeLayout {
  if (data.persons.length === 0) {
    return { nodes: [], edges: data.relationships, width: 400, height: 200 };
  }

  const personById = new Map(data.persons.map((p) => [p.id, p]));
  const generations = new Map<string, number>();

  for (const person of data.persons) {
    generations.set(person.id, 0);
  }

  const parentRels = data.relationships.filter((r) => r.type === "parent");
  const spouseRels = data.relationships.filter((r) => r.type === "spouse");

  // Propagate generations downward from roots.
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

  // Align spouses on the same generation (use the higher row / smaller gen number).
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

  const sortedGens = [...byGeneration.keys()].sort((a, b) => a - b);
  const maxRow = Math.max(...sortedGens.map((g) => byGeneration.get(g)?.length ?? 0), 1);
  const maxRowWidth = maxRow * NODE_WIDTH + Math.max(0, maxRow - 1) * H_GAP;
  const nodes: TreeNodePosition[] = [];

  for (const gen of sortedGens) {
    const row = byGeneration.get(gen) ?? [];
    row.sort((a, b) => displayName(a).localeCompare(displayName(b)));
    const rowWidth = row.length * NODE_WIDTH + Math.max(0, row.length - 1) * H_GAP;
    const startX = PADDING + (maxRowWidth - rowWidth) / 2;
    row.forEach((person, index) => {
      nodes.push({
        person,
        generation: gen,
        x: startX + index * (NODE_WIDTH + H_GAP),
        y: PADDING + gen * (NODE_HEIGHT + V_GAP),
      });
    });
  }

  const maxGen = sortedGens.at(-1) ?? 0;
  const width = Math.max(400, PADDING * 2 + maxRowWidth);
  const height = PADDING * 2 + (maxGen + 1) * NODE_HEIGHT + maxGen * V_GAP;

  return {
    nodes,
    edges: data.relationships.filter(
      (rel) => personById.has(rel.from_person_id) && personById.has(rel.to_person_id),
    ),
    width,
    height,
  };
}

export function displayName(person: Person): string {
  return person.last_name ? `${person.first_name} ${person.last_name}` : person.first_name;
}

export function nodeCenter(node: TreeNodePosition): { x: number; y: number } {
  return { x: node.x + NODE_WIDTH / 2, y: node.y + NODE_HEIGHT / 2 };
}

export function findNode(nodes: TreeNodePosition[], personId: string): TreeNodePosition | undefined {
  return nodes.find((n) => n.person.id === personId);
}

export function edgePath(from: TreeNodePosition, to: TreeNodePosition, type: Relationship["type"]): string {
  const start = nodeCenter(from);
  const end = nodeCenter(to);

  if (type === "spouse") {
    return `M ${start.x} ${start.y} L ${end.x} ${end.y}`;
  }

  const midY = (start.y + end.y) / 2;
  return `M ${start.x} ${start.y + NODE_HEIGHT / 2} C ${start.x} ${midY}, ${end.x} ${midY}, ${end.x} ${end.y - NODE_HEIGHT / 2}`;
}

export const treeMetrics = { NODE_WIDTH, NODE_HEIGHT };
