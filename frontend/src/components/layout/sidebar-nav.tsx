"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { workspaceNavItems } from "@/components/layout/workspace-nav";
import { cn } from "@/lib/utils";

export function SidebarNav() {
  const pathname = usePathname();

  return (
    <nav className="space-y-2">
      {workspaceNavItems.map((item) => {
        const active = pathname === item.href;

        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "flex items-start gap-3 rounded-[1rem] border px-4 py-3 transition-colors",
              active
                ? "border-primary/15 bg-primary/[0.06] text-foreground"
                : "border-transparent text-foreground hover:border-border hover:bg-white/60",
            )}
          >
            <div className={cn("rounded-xl p-2", active ? "bg-primary text-primary-foreground" : "bg-secondary text-primary")}>
              <item.icon className="h-4 w-4" />
            </div>
            <div className="min-w-0">
              <p className="text-sm font-semibold">{item.label}</p>
            </div>
          </Link>
        );
      })}
    </nav>
  );
}
