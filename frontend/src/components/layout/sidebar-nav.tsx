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
              "flex items-start gap-3 rounded-[1.2rem] px-4 py-3 transition-all",
              active
                ? "liquid-button text-white shadow-md"
                : "glass-chip text-foreground hover:-translate-y-0.5 hover:bg-white/70",
            )}
          >
            <div className={cn("rounded-xl p-2", active ? "bg-white/10 text-amber-300" : "bg-secondary text-primary")}>
              <item.icon className="h-4 w-4" />
            </div>
            <div className="min-w-0">
              <p className="text-sm font-semibold">{item.label}</p>
              <p className={cn("mt-1 text-xs leading-5", active ? "text-slate-300" : "text-muted-foreground")}>{item.detail}</p>
            </div>
          </Link>
        );
      })}
    </nav>
  );
}
