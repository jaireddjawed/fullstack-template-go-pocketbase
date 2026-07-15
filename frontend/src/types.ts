import type { AuthUser } from "@shared/types.gen";

// Props gonertia shares with every page (see internal/web/actions.go).
export interface SharedProps {
  auth: AuthUser | null;
  errors: Record<string, string>;
  [key: string]: unknown;
}
