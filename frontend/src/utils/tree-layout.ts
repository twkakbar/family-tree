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
    const household = householdOnRow(person, gen, generations, spouseRels, people);
    for (const member of household) used.add(member.id);
    units.push(household.length > 1 ? orderHousehold(household, parentRels) : household);
  }

  units.sort((a, b) => {
    const ax = unitAnchor(a, parentRels, spouseRels, generations, placedX);
    const bx = unitAnchor(b, parentRels, spouseRels, generations, placedX);
    if (ax !== bx) return ax - bx;
    return displayName(a[0]).localeCompare(displayName(b[0]));
  });

  return units.flat();
}

function householdOnRow(
  person: Person,
  gen: number,
  generations: Map<string, number>,
  spouseRels: Relationship[],
  people: Person[],
): Person[] {
  const onRow = new Map(people.map((candidate) => [candidate.id, candidate]));
  const seen = new Set<string>();
  const queue = [person.id];
  const members: Person[] = [];

  while (queue.length > 0) {
    const id = queue.pop();
    if (!id || seen.has(id)) continue;
    seen.add(id);
    const member = onRow.get(id);
    if (!member) continue;
    members.push(member);
    for (const rel of spouseRels) {
      const otherId =
        rel.from_person_id === id ? rel.to_person_id : rel.to_person_id === id ? rel.from_person_id : undefined;
      if (!otherId || !onRow.has(otherId) || (generations.get(otherId) ?? 0) !== gen) continue;
      queue.push(otherId);
    }
  }

  return members;
}

/** Wives sit beside the bloodline person, split to both sides when there are several. */
function orderHousehold(members: Person[], parentRels: Relationship[]): Person[] {
  const byName = (a: Person, b: Person) => displayName(a).localeCompare(displayName(b));
  const blood = members.filter((person) => hasParent(person.id, parentRels)).sort(byName);
  const marriedIn = members.filter((person) => !hasParent(person.id, parentRels)).sort(byName);
  if (blood.length === 0) return [...members].sort(byName);
  const leftCount = Math.ceil(marriedIn.length / 2);
  return [...marriedIn.slice(0, leftCount), ...blood, ...marriedIn.slice(leftCount)];
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
  const parentXs: number[] = [];
  for (const person of unit) {
    for (const rel of parentRels) {
      if (rel.to_person_id !== person.id) continue;
      const parentX = placedX.get(rel.from_person_id);
      if (parentX !== undefined) parentXs.push(parentX);
    }
  }
  if (parentXs.length === 0) return Number.POSITIVE_INFINITY;
  if (parentXs.length === 1) {
    const onlyParent = parentIdOf(unit, parentRels, placedX);
    return onlyParent ? coupleMidpoint(onlyParent, spouseRels, generations, placedX) ?? parentXs[0] : parentXs[0];
  }
  return parentXs.reduce((sum, value) => sum + value, 0) / parentXs.length;
}

function parentIdOf(
  unit: Person[],
  parentRels: Relationship[],
  placedX: Map<string, number>,
): string | undefined {
  for (const person of unit) {
    for (const rel of parentRels) {
      if (rel.to_person_id === person.id && placedX.has(rel.from_person_id)) return rel.from_person_id;
    }
  }
  return undefined;
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
  if (gen === undefined) return own;
  const spouses = spouseRels.flatMap((edge) => {
    const otherId =
      edge.from_person_id === personId
        ? edge.to_person_id
        : edge.to_person_id === personId
          ? edge.from_person_id
          : undefined;
    if (!otherId || (generations.get(otherId) ?? -1) !== gen) return [];
    const other = placedX.get(otherId);
    return other === undefined ? [] : [other];
  });
  if (spouses.length !== 1) return own;
  return (own + spouses[0]) / 2;
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
    const gap = right.x - (left.x + NODE_WIDTH);
    const sameRow = Math.abs(left.y - right.y) < 2;
    if (!sameRow || gap <= H_GAP + 8) {
      return `M ${left.x + NODE_WIDTH} ${y} L ${right.x} ${y}`;
    }
    const arcY = Math.min(left.y, right.y) - 28;
    const midX = (left.x + NODE_WIDTH + right.x) / 2;
    return `M ${left.x + NODE_WIDTH} ${y} Q ${midX} ${arcY}, ${right.x} ${y}`;
  }

  const midY = (start.y + end.y) / 2;
  return `M ${start.x} ${start.y + NODE_HEIGHT / 2} C ${start.x} ${midY}, ${end.x} ${midY}, ${end.x} ${end.y - NODE_HEIGHT / 2}`;
}

export const treeMetrics = { NODE_WIDTH, NODE_HEIGHT };
