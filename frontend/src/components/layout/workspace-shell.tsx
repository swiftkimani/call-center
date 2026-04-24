"use client";

import type { ReactNode } from "react";
import { useMemo, useState } from "react";
import { Bell, Headphones, Search, ShieldCheck } from "lucide-react";

import { SidebarNav } from "@/components/layout/sidebar-nav";
import { TopNav } from "@/components/layout/top-nav";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const quickSearch = [
  { label: "Agent desk", href: "/agent" },
  { label: "Supervisor board", href: "/supervisor" },
  { label: "Campaign planner", href: "/campaigns" },
  { label: "Reports", href: "/reports" },
  { label: "Settings", href: "/settings" },
];

export function WorkspaceShell({ children }: { children: ReactNode }) {
  const [query, setQuery] = useState("");
  const [showAlerts, setShowAlerts] = useState(false);

  const matches = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    if (!normalized) return [];

    return quickSearch.filter((item) => item.label.toLowerCase().includes(normalized));
  }, [query]);

  return (
    <main className="mx-auto min-h-screen max-w-[1600px] px-4 py-4 md:px-6">
      <div className="grid min-h-[calc(100vh-2rem)] gap-4 lg:grid-cols-[280px_1fr]">
        <aside className="glass-panel rounded-[2rem] p-4">
          <div className="glass-chip flex items-center justify-between rounded-[1.6rem] px-4 py-4">
            <div className="flex items-center gap-3">
              <div className="liquid-button rounded-2xl p-3 text-amber-300">
                <Headphones className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm font-semibold text-foreground">Call Center</p>
                <p className="text-xs text-muted-foreground">Live operations</p>
              </div>
            </div>
            <div className="rounded-full bg-emerald-500/15 px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.18em] text-emerald-700">
              Online
            </div>
          </div>

          <div className="mt-4">
            <SidebarNav />
          </div>

          <div className="glass-chip mt-4 rounded-[1.4rem] p-4">
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-primary" />
              <p className="text-sm font-semibold">System status</p>
            </div>
            <p className="mt-2 text-sm leading-6 text-muted-foreground">
              Frontend remains available if API or worker services are degraded.
            </p>
          </div>
        </aside>

        <section className="glass-panel rounded-[2rem]">
          <header className="border-b border-border/60 px-5 py-4 md:px-6">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
              <div className="space-y-4">
                <div>
                  <h1 className="text-2xl font-semibold tracking-tight text-foreground">Workspace</h1>
                  <p className="mt-1 text-sm leading-6 text-muted-foreground">Routes are grouped by role so the navigation stays shallow.</p>
                </div>
                <TopNav />
              </div>

              <div className="w-full max-w-[420px] space-y-3">
                <div className="flex gap-3">
                  <div className="relative flex-1">
                    <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input value={query} onChange={(event) => setQuery(event.target.value)} className="pl-9" placeholder="Find a route" />
                  </div>
                  <Button variant="glass" size="default" type="button" onClick={() => setShowAlerts((value) => !value)} aria-pressed={showAlerts}>
                    <Bell className="h-4 w-4 text-primary" />
                  </Button>
                </div>

                {query && (
                  <div className="glass-chip rounded-[1.2rem] px-3 py-3 text-sm text-muted-foreground">
                    {matches.length > 0 ? (
                      <div className="flex flex-wrap gap-2">
                        {matches.map((item) => (
                          <a key={item.href} href={item.href} className="rounded-full bg-white/70 px-3 py-1.5 font-semibold text-foreground transition hover:bg-white">
                            {item.label}
                          </a>
                        ))}
                      </div>
                    ) : (
                      <p>No matching route.</p>
                    )}
                  </div>
                )}

                {showAlerts && (
                  <div className="glass-chip rounded-[1.2rem] px-4 py-3 text-sm text-muted-foreground">
                    <p className="font-semibold text-foreground">Alerts</p>
                    <p className="mt-1">No active incidents. Queue spikes and dropped websocket sessions will appear here.</p>
                  </div>
                )}
              </div>
            </div>
          </header>

          <div className="px-5 py-5 md:px-6">{children}</div>
        </section>
      </div>
    </main>
  );
}
