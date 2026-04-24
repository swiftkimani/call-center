import type { LucideIcon } from "lucide-react";
import { BarChart3, Headphones, Layers3, RadioTower, Settings2 } from "lucide-react";

export type WorkspaceNavItem = {
  href: string;
  label: string;
  detail: string;
  shortLabel: string;
  icon: LucideIcon;
};

export const workspaceNavItems: WorkspaceNavItem[] = [
  { href: "/agent", label: "Agent Workspace", shortLabel: "Agent", detail: "Live call, notes, softphone", icon: Headphones },
  { href: "/supervisor", label: "Supervisor", shortLabel: "Supervisor", detail: "Queues, agents, interventions", icon: RadioTower },
  { href: "/campaigns", label: "Campaigns", shortLabel: "Campaigns", detail: "Outbound programs and imports", icon: Layers3 },
  { href: "/reports", label: "Reports", shortLabel: "Reports", detail: "KPIs, exports, performance", icon: BarChart3 },
  { href: "/settings", label: "Settings", shortLabel: "Settings", detail: "Queues, routing, platform", icon: Settings2 },
];
