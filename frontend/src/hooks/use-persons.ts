"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { queryKeys } from "@/lib/query-keys";
import { personsService } from "@/services/persons.service";
import type {
  CreatePersonInput,
  CreateRelationshipInput,
  UpdatePersonInput,
} from "@/types/person";

export function useFamilyTree(familyId: string) {
  return useQuery({
    queryKey: queryKeys.families.tree(familyId),
    queryFn: () => personsService.getTree(familyId),
    enabled: Boolean(familyId),
  });
}

export function useFamilyPersons(familyId: string, search = "", page = 1) {
  return useQuery({
    queryKey: queryKeys.families.persons(familyId, search, page),
    queryFn: () => personsService.listByFamily(familyId, search, page),
    enabled: Boolean(familyId),
  });
}

export function usePerson(personId: string) {
  return useQuery({
    queryKey: queryKeys.persons.detail(personId),
    queryFn: () => personsService.getById(personId),
    enabled: Boolean(personId),
  });
}

export function usePersonAncestors(personId: string, depth = 5) {
  return useQuery({
    queryKey: queryKeys.persons.ancestors(personId, depth),
    queryFn: () => personsService.getAncestors(personId, depth),
    enabled: Boolean(personId),
  });
}

export function usePersonDescendants(personId: string, depth = 5) {
  return useQuery({
    queryKey: queryKeys.persons.descendants(personId, depth),
    queryFn: () => personsService.getDescendants(personId, depth),
    enabled: Boolean(personId),
  });
}

function invalidateFamilyPersonData(queryClient: ReturnType<typeof useQueryClient>, familyId: string) {
  queryClient.invalidateQueries({ queryKey: queryKeys.families.tree(familyId) });
  queryClient.invalidateQueries({ queryKey: ["families", "persons", familyId] });
  queryClient.invalidateQueries({ queryKey: queryKeys.families.detail(familyId) });
  queryClient.invalidateQueries({ queryKey: queryKeys.families.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.public.tree });
  queryClient.invalidateQueries({ queryKey: queryKeys.public.meta });
}

export function useCreatePerson(familyId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreatePersonInput) => personsService.create(familyId, input),
    onSuccess: () => invalidateFamilyPersonData(queryClient, familyId),
  });
}

export function useUpdatePerson(familyId: string, personId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: UpdatePersonInput) => personsService.update(personId, input),
    onSuccess: () => {
      invalidateFamilyPersonData(queryClient, familyId);
      queryClient.invalidateQueries({ queryKey: queryKeys.persons.detail(personId) });
    },
  });
}

export function useDeletePerson(familyId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (personId: string) => personsService.remove(personId),
    onSuccess: () => invalidateFamilyPersonData(queryClient, familyId),
  });
}

export function useUploadPersonPhoto(familyId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ personId, file }: { personId: string; file: File }) =>
      personsService.uploadPhoto(personId, file),
    onSuccess: (person) => {
      invalidateFamilyPersonData(queryClient, familyId);
      queryClient.invalidateQueries({ queryKey: queryKeys.persons.detail(person.id) });
    },
  });
}

export function useCreateRelationship(familyId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateRelationshipInput) =>
      personsService.createRelationship(familyId, input),
    onSuccess: () => invalidateFamilyPersonData(queryClient, familyId),
  });
}

export function useDeleteRelationship(familyId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (relationshipId: string) => personsService.removeRelationship(relationshipId),
    onSuccess: () => invalidateFamilyPersonData(queryClient, familyId),
  });
}
