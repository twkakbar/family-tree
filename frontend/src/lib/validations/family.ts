import { z } from "zod";

export const createFamilySchema = z.object({
  name: z.string().trim().min(2, "Name must be at least 2 characters").max(100),
  description: z.string().trim().max(1000).optional().or(z.literal("")),
});

export const updateFamilySchema = createFamilySchema;

export const inviteCollaboratorSchema = z.object({
  email: z.string().trim().email("Enter a valid email address"),
  role: z.enum(["editor", "viewer"]),
});

export type CreateFamilyFormValues = z.infer<typeof createFamilySchema>;
export type UpdateFamilyFormValues = z.infer<typeof updateFamilySchema>;
export type InviteCollaboratorFormValues = z.infer<typeof inviteCollaboratorSchema>;
