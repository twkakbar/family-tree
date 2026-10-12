import type { Relationship, TreeNodePosition } from "@/types/person";
import { nodeCenter, treeMetrics } from "@/utils/tree-layout";

const { NODE_WIDTH, NODE_HEIGHT } = treeMetrics;

export interface KinshipEdge {
  from_person_id: string;
  to_person_id: string;
  type: Relationship["type"];
}

export interface SpousePair {
  a: string;
  b: string;
}

export interface DescentGroup {
  parentIds: string[];
  childIds: string[];
}

/** Couples, plus sibling groups that share the same parents. */
export function kinshipLayout(edges: KinshipEdge[]): {
  spouses: SpousePair[];
  descents: DescentGroup[];
} {
  const spouses = edges
    .filter((edge) => edge.type === "spouse")
    .map((edge) => ({ a: edge.from_person_id, b: edge.to_person_id }));

  const parentsOf = new Map<string, string[]>();
  for (const edge of edges) {
    if (edge.type !== "parent") continue;
    const list = parentsOf.get(edge.to_person_id) ?? [];
    list.push(edge.from_person_id);
    parentsOf.set(edge.to_person_id, list);
  }

  const groups = new Map<string, DescentGroup>();
  for (const [childId, parentIds] of parentsOf) {
    const unique = [...new Set(parentIds)].sort();
    if (unique.length === 0) continue;
    const key = unique.join("|");
    const group = groups.get(key) ?? { parentIds: unique, childIds: [] };
    group.childIds.push(childId);
    groups.set(key, group);
  }

  return { spouses, descents: [...groups.values()] };
}

export interface PedigreePath {
  kind: "spouse" | "descent";
  d: string;
}

/** Solid marriage lines, and one shared branch for children of the same parents. */
export function studioPedigree(nodes: TreeNodePosition[], edges: KinshipEdge[]): PedigreePath[] {
  const byId = new Map(nodes.map((node) => [node.person.id, node]));
  const { spouses, descents } = kinshipLayout(edges);
  const paths: PedigreePath[] = [];

  for (const pair of spouses) {
    const left = byId.get(pair.a);
    const right = byId.get(pair.b);
    if (!left || !right) continue;
    paths.push({ kind: "spouse", d: spousePath(left, right, nodes) });
  }

  for (const group of descents) {
    const parents = group.parentIds
      .map((id) => byId.get(id))
      .filter((node): node is TreeNodePosition => Boolean(node));
    const children = group.childIds
      .map((id) => byId.get(id))
      .filter((node): node is TreeNodePosition => Boolean(node));
    if (parents.length === 0 || children.length === 0) continue;
    paths.push({ kind: "descent", d: descentPath(parents, children, nodes) });
  }

  return paths;
}

function spousePath(a: TreeNodePosition, b: TreeNodePosition, nodes: TreeNodePosition[]): string {
  const left = a.x <= b.x ? a : b;
  const right = a.x <= b.x ? b : a;
  const y = (left.y + NODE_HEIGHT + right.y + NODE_HEIGHT) / 2;
  const x1 = left.x + NODE_WIDTH;
  const x2 = right.x;
  const sameRow = Math.abs(left.y - right.y) < 2;
  const blocked =
    sameRow &&
    nodes.some((node) => {
      if (node.person.id === left.person.id || node.person.id === right.person.id) return false;
      if (Math.abs(node.y - left.y) > 2) return false;
      const center = node.x + NODE_WIDTH / 2;
      return center > x1 && center < x2;
    });
  if (!sameRow || !blocked) return `M ${x1} ${y} L ${x2} ${y}`;
  const arcY = Math.min(left.y, right.y) - 28;
  return `M ${x1} ${y} Q ${(x1 + x2) / 2} ${arcY}, ${x2} ${y}`;
}

function descentPath(
  parents: TreeNodePosition[],
  children: TreeNodePosition[],
  nodes: TreeNodePosition[],
): string {
  const parentBottom = Math.max(...parents.map((parent) => parent.y + NODE_HEIGHT));
  const childTop = Math.min(...children.map((child) => child.y));
  const gap = Math.max(36, childTop - parentBottom);
  const barY = parentBottom + gap * 0.45;

  let dropX = nodeCenter(parents[0]).x;
  let stemTop = parents[0].y + NODE_HEIGHT;

  if (parents.length >= 2) {
    const ordered = [...parents].sort((a, b) => a.x - b.x);
    const left = ordered[0];
    const right = ordered[ordered.length - 1];
    const x1 = left.x + NODE_WIDTH;
    const x2 = right.x;
    const between = nodes
      .filter((node) => {
        if (parents.some((parent) => parent.person.id === node.person.id)) return false;
        if (Math.abs(node.y - left.y) > 2) return false;
        const center = node.x + NODE_WIDTH / 2;
        return center > x1 && center < x2;
      })
      .sort((a, b) => a.x - b.x);
    const sameRow = Math.abs(left.y - right.y) < 2;
    if (sameRow && between.length === 0) {
      dropX = (x1 + x2) / 2;
      stemTop = (left.y + NODE_HEIGHT + right.y + NODE_HEIGHT) / 2;
    } else if (between.length > 0) {
      dropX = (x1 + between[0].x) / 2;
      stemTop = parentBottom;
    } else {
      dropX = parents.reduce((sum, parent) => sum + nodeCenter(parent).x, 0) / parents.length;
      stemTop = parentBottom;
    }
  }

  const childXs = children.map((child) => nodeCenter(child).x);
  const barLeft = Math.min(dropX, ...childXs);
  const barRight = Math.max(dropX, ...childXs);
  const parts = [`M ${dropX} ${stemTop} L ${dropX} ${barY}`];
  if (barRight - barLeft > 1) parts.push(`M ${barLeft} ${barY} L ${barRight} ${barY}`);
  for (const child of children) {
    const x = nodeCenter(child).x;
    parts.push(`M ${x} ${barY} L ${x} ${child.y}`);
  }
  return parts.join(" ");
}
