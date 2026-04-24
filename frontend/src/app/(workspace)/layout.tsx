import type { ReactNode } from "react";
import { redirect } from "next/navigation";

import { WorkspaceShell } from "@/components/layout/workspace-shell";
import { hasSession } from "@/lib/auth/session";

export default async function WorkspaceLayout({ children }: { children: ReactNode }) {
  if (!(await hasSession())) {
    redirect("/");
  }

  return <WorkspaceShell>{children}</WorkspaceShell>;
}
