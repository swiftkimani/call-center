"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { workspaceNavItems } from "@/components/layout/workspace-nav";
import { cn } from "@/lib/utils";

export function TopNav() {
  const pathname = usePathname();

  return (
    <nav aria-label="Primary" className="overflow-x-auto">
      <div className="inline-flex min-w-full gap-1 rounded-[1rem] bg-white/72 p-1">
        {workspaceNavItems.map((item) => {
          const active = pathname === item.href;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "inline-flex items-center gap-2 rounded-[0.8rem] px-4 py-2.5 text-sm font-semibold whitespace-nowrap transition-colors",
                active ? "bg-primary text-white" : "text-muted-foreground hover:bg-secondary hover:text-foreground",
              )}
            >
              <item.icon className={cn("h-4 w-4", active ? "text-white" : "text-muted-foreground")} />
              <span className={active ? "text-white" : undefined}>{item.shortLabel}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
