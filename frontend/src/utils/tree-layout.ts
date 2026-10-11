import type { FamilyTreeData, Person, Relationship, TreeLayout, TreeNodePosition } from "@/types/person";

const NODE_WIDTH = 152;
const NODE_HEIGHT = 76;
const H_GAP = 48;
const V_GAP = 96;
const PADDING = 40;

/** Generation 0 is the oldest row. Spouses share a row; children sit strictly below parents. */
export function assignGenerations(data: FamilyTreeData): Map<string, number> {
  const generations = new Map<string, number>();
  for (const person of data.persons) generations.set(person.id, 0);

  const parentRels = data.relationships.filter((rel) => rel.type === "parent");
  const spouseRels = data.relationships.filter((rel) => rel.type === "spouse");

  const propagateParents = () => {
    let changed = true;
    let guard = 0;
    while (changed && guard < data.persons.length + 2) {
      changed = false;
      guard += 1;
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
  };

  propagateParents();

  let guard = 0;
  let moved = true;
  while (moved && guard < data.persons.length + 2) {
    moved = false;
    guard += 1;
    for (const rel of spouseRels) {
      const a = generations.get(rel.from_person_id) ?? 0;
      const b = generations.get(rel.to_person_id) ?? 0;
      // Keep the bloodline row. A spouse with no parents must come down
      // beside their partner, not pull the partner up to the top.
      const row = Math.max(a, b);
      if (a !== row || b !== row) {
        generations.set(rel.from_person_id, row);
        generations.set(rel.to_person_id, row);
        moved = true;
      }
    }
    if (moved) propagateParents();
  }

  const min = Math.min(...generations.values());
  if (min > 0) {
    for (const [id, gen] of generations) generations.set(id, gen - min);
  }
  return generations;
}

/** People in each generation, spouses adjacent and rows ordered under their parents. */
export function peopleByGeneration(data: FamilyTreeData): Person[][] {
  const generations = assignGenerations(data);
  const parentRels = data.relationships.filter((rel) => rel.type === "parent");
  const spouseRels = data.relationships.filter((rel) => rel.type === "spouse");
  const byGeneration = new Map<number, Person[]>();

  for (const person of data.persons) {
    const gen = generations.get(person.id) ?? 0;
    const bucket = byGeneration.get(gen) ?? [];
    bucket.push(person);
    byGeneration.set(gen, bucket);
  }

  const sortedGens = [...byGeneration.keys()].sort((a, b) => a - b);
  const placedX = new Map<string, number>();
  const rows: Person[][] = [];

  for (const gen of sortedGens) {
    const people = byGeneration.get(gen) ?? [];
    const row = orderGeneration(people, gen, generations, parentRels, spouseRels, placedX);
    row.forEach((person, index) => placedX.set(person.id, index));
    rows.push(row);
  }

  return rows;
}

function orderGeneration(
  people: Person[],
  gen: number,
  generations: Map<string, number>,
  parentRels: Relationship[],
  spouseRels: Relationship[],
  placedX: Map<string, number>,
): Person[] {
  const used = new Set<string>();
  const units: Person[][] = [];

  for (const person of people) {
    if (used.has(person.id)) continue;
    const partner = spouseOnRow(person, gen, generations, spouseRels, people);
    if (partner) {
      used.add(person.id);
      used.add(partner.id);
      units.push(orderCouple(person, partner, parentRels));
    } else {
      used.add(person.id);
      units.push([person]);
    }
  }

  units.sort((a, b) => {
    const ax = unitAnchor(a, parentRels, spouseRels, generations, placedX);
    const bx = unitAnchor(b, parentRels, spouseRels, generations, placedX);
    if (ax !== bx) return ax - bx;
    return displayName(a[0]).localeCompare(displayName(b[0]));
  });

  return units.flat();
}

function spouseOnRow(
  person: Person,
  gen: number,
  generations: Map<string, number>,
  spouseRels: Relationship[],
  people: Person[],
): Person | undefined {
  const rel = spouseRels.find(
    (edge) => edge.from_person_id === person.id || edge.to_person_id === person.id,
  );
  if (!rel) return undefined;
  const otherId = rel.from_person_id === person.id ? rel.to_person_id : rel.from_person_id;
  if ((generations.get(otherId) ?? 0) !== gen) return undefined;
  return people.find((candidate) => candidate.id === otherId);
}

function orderCouple(a: Person, b: Person, parentRels: Relationship[]): Person[] {
  const aBlood = hasParent(a.id, parentRels);
  const bBlood = hasParent(b.id, parentRels);
  if (aBlood !== bBlood) return aBlood ? [b, a] : [a, b];
  return displayName(a).localeCompare(displayName(b)) <= 0 ? [a, b] : [b, a];
}

function hasParent(personId: string, parentRels: Relationship[]): boolean {
  return parentRels.some((rel) => rel.to_person_id === personId);
}

function unitAnchor(
  unit: Person[],
  parentRels: Relationship[],
  spouseRels: Relationship[],
  generations: Map<string, number>,
  placedX: Map<string, number>,
): number {
  const anchors: number[] = [];
  for (const person of unit) {
    for (const rel of parentRels) {
      if (rel.to_person_id !== person.id) continue;
      const parentX = coupleMidpoint(rel.from_person_id, spouseRels, generations, placedX);
      if (parentX !== undefined) anchors.push(parentX);
    }
  }
  if (anchors.length === 0) return Number.POSITIVE_INFINITY;
  return anchors.reduce((sum, value) => sum + value, 0) / anchors.length;
}

function coupleMidpoint(
  personId: string,
  spouseRels: Relationship[],
  generations: Map<string, number>,
  placedX: Map<string, number>,
): number | undefined {
  const own = placedX.get(personId);
  if (own === undefined) return undefined;
  const gen = generations.get(personId);
  const rel = spouseRels.find(
    (edge) => edge.from_person_id === personId || edge.to_person_id === personId,
  );
  if (!rel || gen === undefined) return own;
  const otherId = rel.from_person_id === personId ? rel.to_person_id : rel.from_person_id;
  if ((generations.get(otherId) ?? -1) !== gen) return own;
  const other = placedX.get(otherId);
  if (other === undefined) return own;
  return (own + other) / 2;
}

/** Assigns generation rows and x positions for a simple top-down tree layout. */
export function layoutFamilyTree(data: FamilyTreeData): TreeLayout {
  if (data.persons.length === 0) {
    return { nodes: [], edges: data.relationships, width: 400, height: 200 };
  }

  const personById = new Map(data.persons.map((p) => [p.id, p]));
  const rows = peopleByGeneration(data);
  const maxRow = Math.max(...rows.map((row) => row.length), 1);
  const maxRowWidth = maxRow * NODE_WIDTH + Math.max(0, maxRow - 1) * H_GAP;
  const nodes: TreeNodePosition[] = [];

  rows.forEach((row, gen) => {
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
  });

  const maxGen = Math.max(0, rows.length - 1);
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
    const left = start.x <= end.x ? from : to;
    const right = start.x <= end.x ? to : from;
    const y = (nodeCenter(left).y + nodeCenter(right).y) / 2;
    return `M ${left.x + NODE_WIDTH} ${y} L ${right.x} ${y}`;
  }

  const midY = (start.y + end.y) / 2;
  return `M ${start.x} ${start.y + NODE_HEIGHT / 2} C ${start.x} ${midY}, ${end.x} ${midY}, ${end.x} ${end.y - NODE_HEIGHT / 2}`;
}

export const treeMetrics = { NODE_WIDTH, NODE_HEIGHT };
