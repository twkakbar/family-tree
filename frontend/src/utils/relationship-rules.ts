import type { Person, Relationship } from "@/types/person";
import { displayName } from "@/utils/tree-layout";

export type RelationshipLinkChoice = "parent-of-target" | "child-of-target" | "spouse";

function parentsOf(personId: string, rels: Relationship[]): string[] {
  return rels
    .filter((rel) => rel.type === "parent" && rel.to_person_id === personId)
    .map((rel) => rel.from_person_id);
}

function hasSpouse(personId: string, rels: Relationship[]): boolean {
  return rels.some(
    (rel) =>
      rel.type === "spouse" &&
      (rel.from_person_id === personId || rel.to_person_id === personId),
  );
}

/** A husband may take another wife. Anyone else keeps a single spouse. */
export function canHaveAnotherSpouse(person: Person, rels: Relationship[]): boolean {
  if (person.gender === "male") return true;
  return !hasSpouse(person.id, rels);
}

function isAncestor(personId: string, maybeAncestor: string, rels: Relationship[]): boolean {
  const seen = new Set<string>();
  const stack = parentsOf(personId, rels);
  while (stack.length > 0) {
    const id = stack.pop();
    if (!id || seen.has(id)) continue;
    if (id === maybeAncestor) return true;
    seen.add(id);
    stack.push(...parentsOf(id, rels));
  }
  return false;
}

function shareParent(a: string, b: string, rels: Relationship[]): boolean {
  const parents = new Set(parentsOf(a, rels));
  return parentsOf(b, rels).some((id) => parents.has(id));
}

function bornAfter(parent: Person, child: Person): boolean {
  if (!parent.birth_date || !child.birth_date) return false;
  return parent.birth_date > child.birth_date;
}

/** Choices that would not break the lineage. Mirrors the server rules. */
export function allowedLinkChoices(
  source: Person,
  target: Person,
  rels: Relationship[],
): RelationshipLinkChoice[] {
  const alreadySpouse = rels.some(
    (rel) =>
      rel.type === "spouse" &&
      ((rel.from_person_id === source.id && rel.to_person_id === target.id) ||
        (rel.from_person_id === target.id && rel.to_person_id === source.id)),
  );
  const alreadyParent = (fromId: string, toId: string) =>
    rels.some((rel) => rel.type === "parent" && rel.from_person_id === fromId && rel.to_person_id === toId);

  const canBeParent = (parent: Person, child: Person) => {
    if (alreadyParent(parent.id, child.id) || alreadyParent(child.id, parent.id)) return false;
    if (alreadySpouse) return false;
    if (parentsOf(child.id, rels).length >= 2) return false;
    if (isAncestor(parent.id, child.id, rels)) return false;
    if (bornAfter(parent, child)) return false;
    return true;
  };

  const choices: RelationshipLinkChoice[] = [];
  if (canBeParent(source, target)) choices.push("parent-of-target");
  if (canBeParent(target, source)) choices.push("child-of-target");

  const blood =
    isAncestor(source.id, target.id, rels) || isAncestor(target.id, source.id, rels);
  if (
    !alreadySpouse &&
    !blood &&
    !shareParent(source.id, target.id, rels) &&
    canHaveAnotherSpouse(source, rels) &&
    canHaveAnotherSpouse(target, rels)
  ) {
    choices.push("spouse");
  }
  return choices;
}

export function relationshipLabel(rel: Relationship, personId: string, people: Person[]): string {
  const otherId = rel.from_person_id === personId ? rel.to_person_id : rel.from_person_id;
  const other = people.find((person) => person.id === otherId);
  const name = other ? displayName(other) : "another member";
  if (rel.type === "spouse") return `Spouse of ${name}`;
  if (rel.from_person_id === personId) return `Parent of ${name}`;
  return `Child of ${name}`;
}
