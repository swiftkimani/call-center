"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { workspaceNavItems } from "@/components/layout/workspace-nav";
import { cn } from "@/lib/utils";

export function TopNav() {
  const pathname = usePathname();

  return (
    <nav aria-label="Primary" className="flex flex-wrap gap-2">
      {workspaceNavItems.map((item) => {
        const active = pathname === item.href;

        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              "glass-chip inline-flex items-center gap-2 rounded-full px-3 py-2 text-sm font-semibold transition-all",
              active
                ? "liquid-button text-primary-foreground shadow-sm"
                : "text-muted-foreground hover:-translate-y-0.5 hover:text-foreground",
            )}
          >
            <item.icon className={cn("h-4 w-4", active ? "text-amber-300" : "text-primary")} />
            <span>{item.shortLabel}</span>
          </Link>
        );
      })}
    </nav>
  );
}
