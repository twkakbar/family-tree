import { z } from "zod";

export const personSchema = z.object({
  first_name: z.string().trim().min(1, "Name is required").max(100),
  last_name: z.string().trim().max(100).optional().or(z.literal("")),
  gender: z.enum(["male", "female", "other", "unknown"]),
  birth_date: z
    .string()
    .trim()
    .regex(/^\d{4}-\d{2}-\d{2}$/, "Use YYYY-MM-DD")
    .optional()
    .or(z.literal("")),
  death_date: z
    .string()
    .trim()
    .regex(/^\d{4}-\d{2}-\d{2}$/, "Use YYYY-MM-DD")
    .optional()
    .or(z.literal("")),
  bio: z.string().trim().max(5000).optional().or(z.literal("")),
});

export const relationshipSchema = z
  .object({
    from_person_id: z.string().uuid("Select the first person"),
    to_person_id: z.string().uuid("Select the second person"),
    type: z.enum(["parent", "spouse"]),
  })
  .refine((data) => data.from_person_id !== data.to_person_id, {
    message: "Choose two different people",
    path: ["to_person_id"],
  });

export type PersonFormValues = z.infer<typeof personSchema>;
export type RelationshipFormValues = z.infer<typeof relationshipSchema>;
