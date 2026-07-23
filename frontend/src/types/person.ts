export type Gender = "male" | "female" | "other" | "unknown";
export type RelationshipType = "parent" | "spouse";

export interface Person {
  id: string;
  family_id: string;
  first_name: string;
  last_name: string;
  gender: Gender;
  birth_date: string | null;
  death_date: string | null;
  bio: string;
  photo_url: string;
  created_at: string;
  updated_at: string;
}

export interface Relationship {
  id: string;
  family_id: string;
  from_person_id: string;
  to_person_id: string;
  type: RelationshipType;
  created_at: string;
}

export interface PersonDetail extends Person {
  relationships: Relationship[];
}

export interface PersonWithDepth extends Person {
  depth: number;
}

export interface FamilyTreeData {
  persons: Person[];
  relationships: Relationship[];
}

export interface CreatePersonInput {
  first_name: string;
  last_name?: string;
  gender: Gender;
  birth_date?: string;
  death_date?: string;
  bio?: string;
}

export type UpdatePersonInput = CreatePersonInput;

export interface CreateRelationshipInput {
  from_person_id: string;
  to_person_id: string;
  type: RelationshipType;
}

export interface PersonListResult {
  items: Person[];
  meta: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

export interface TreeNodePosition {
  person: Person;
  x: number;
  y: number;
  generation: number;
}

export interface TreeLayout {
  nodes: TreeNodePosition[];
  edges: Relationship[];
  width: number;
  height: number;
}
