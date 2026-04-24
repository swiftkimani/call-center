"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import { useMemo, useState } from "react";
import { Bell, Menu, Search, ShieldCheck, X } from "lucide-react";

import { BrandLogo } from "@/components/layout/brand-logo";
import { SidebarNav } from "@/components/layout/sidebar-nav";
import { TopNav } from "@/components/layout/top-nav";
import { workspaceNavItems } from "@/components/layout/workspace-nav";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function WorkspaceShell({ children }: { children: ReactNode }) {
  const [query, setQuery] = useState("");
  const [showAlerts, setShowAlerts] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const matches = useMemo(() => {
    const normalized = query.trim().toLowerCase();
    if (!normalized) return [];

    return workspaceNavItems.filter((item) => item.label.toLowerCase().includes(normalized) || item.shortLabel.toLowerCase().includes(normalized));
  }, [query]);

  return (
    <main className="mx-auto min-h-screen max-w-[1600px] px-4 py-4 md:px-6">
      <div className="grid min-h-[calc(100vh-2rem)] gap-5 lg:grid-cols-[240px_1fr]">
        <aside className="glass-panel hidden rounded-[1.75rem] p-4 lg:block">
          <div className="px-2 py-1">
            <BrandLogo />
          </div>

          <div className="mt-5">
            <SidebarNav />
          </div>

          <div className="mt-6 rounded-[1.2rem] border border-border bg-white/50 p-4">
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-primary" />
              <p className="text-sm font-semibold">System status</p>
            </div>
            <p className="mt-2 text-sm leading-6 text-muted-foreground">Frontend remains available if API or worker services are degraded.</p>
          </div>
        </aside>

        <section className="glass-panel rounded-[1.75rem]">
          <header className="px-4 py-4 md:px-6 md:py-5">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
              <div className="flex min-w-0 items-center gap-3">
                <div className="flex items-center gap-3">
                  <Button variant="outline" size="sm" type="button" className="lg:hidden" onClick={() => setMobileMenuOpen(true)} aria-label="Open menu">
                    <Menu className="h-4 w-4" />
                  </Button>
                </div>

                <div className="hidden min-w-0 flex-1 xl:block">
                  <TopNav />
                </div>
              </div>

              <div className="flex items-center gap-2">
                <div className="hidden min-w-[280px] md:block">
                  <div className="relative">
                    <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input value={query} onChange={(event) => setQuery(event.target.value)} className="pl-9" placeholder="Find a route" />
                  </div>
                </div>
                <Button variant="outline" size="sm" type="button" onClick={() => setShowAlerts((value) => !value)} aria-pressed={showAlerts} aria-label="Toggle alerts">
                  <Bell className="h-4 w-4 text-primary" />
                </Button>
              </div>

              <div className="md:hidden">
                {query && (
                  <div className="rounded-[1rem] border border-border bg-white/80 px-3 py-3 text-sm text-muted-foreground">
                    {matches.length > 0 ? (
                      <div className="flex flex-wrap gap-2">
                        {matches.map((item) => (
                          <Link key={item.href} href={item.href} className="rounded-full bg-secondary px-3 py-1.5 font-semibold text-foreground transition hover:bg-secondary/75">
                            {item.shortLabel}
                          </Link>
                        ))}
                      </div>
                    ) : (
                      <p>No matching route.</p>
                    )}
                  </div>
                )}

                {showAlerts && (
                  <div className="mt-3 rounded-[1rem] border border-border bg-white/80 px-4 py-3 text-sm text-muted-foreground">
                    <p className="font-semibold text-foreground">Alerts</p>
                    <p className="mt-1">No active incidents.</p>
                  </div>
                )}
              </div>
            </div>

            <div className="mt-4 hidden md:block">
              {query && (
                <div className="rounded-[1rem] border border-border bg-white/80 px-3 py-3 text-sm text-muted-foreground">
                  {matches.length > 0 ? (
                    <div className="flex flex-wrap gap-2">
                      {matches.map((item) => (
                        <Link key={item.href} href={item.href} className="rounded-full bg-secondary px-3 py-1.5 font-semibold text-foreground transition hover:bg-secondary/75">
                          {item.shortLabel}
                        </Link>
                      ))}
                    </div>
                  ) : (
                    <p>No matching route.</p>
                  )}
                </div>
              )}

              {showAlerts && (
                <div className="rounded-[1rem] border border-border bg-white/80 px-4 py-3 text-sm text-muted-foreground">
                  <p className="font-semibold text-foreground">Alerts</p>
                  <p className="mt-1">No active incidents.</p>
                </div>
              )}
            </div>
          </header>

          <div className="border-t border-border/60 px-4 py-5 md:px-6 md:py-6">{children}</div>
        </section>
      </div>

      {mobileMenuOpen && (
        <div className="fixed inset-0 z-50 bg-black/30 lg:hidden">
          <div className="absolute inset-y-0 left-0 w-[88vw] max-w-[340px] border-r border-border bg-[rgba(250,247,242,0.98)] p-4 shadow-2xl backdrop-blur">
            <div className="flex items-center justify-between">
              <BrandLogo />
              <Button variant="outline" size="sm" type="button" onClick={() => setMobileMenuOpen(false)} aria-label="Close menu">
                <X className="h-4 w-4" />
              </Button>
            </div>

            <div className="mt-4">
              <div className="relative">
                <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input value={query} onChange={(event) => setQuery(event.target.value)} className="pl-9" placeholder="Find a route" />
              </div>
            </div>

            <div className="mt-5">
              <SidebarNav />
            </div>
          </div>
          <button type="button" className="absolute inset-0 -z-10" onClick={() => setMobileMenuOpen(false)} aria-label="Close menu overlay" />
        </div>
      )}
    </main>
  );
}
