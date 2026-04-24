import type { ReactNode } from "react";
import { Activity, BarChart3, PhoneCall, ShieldCheck, Waves } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

const pillars = [
  {
    icon: ShieldCheck,
    title: "Operational trust",
    text: "A UI that stays useful during outages, retries, and partial backend availability.",
  },
  {
    icon: Waves,
    title: "Realtime-ready",
    text: "HTTP for workflows, WebSockets for supervision, and brokers for async execution.",
  },
  {
    icon: BarChart3,
    title: "Scalable by slice",
    text: "Frontend features and backend modules can evolve independently without becoming tangled.",
  },
];

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <main className="mx-auto flex min-h-screen max-w-7xl flex-col px-6 py-8 md:px-10">
      <section className="overflow-hidden rounded-[2rem] border border-white/70 bg-white/70 shadow-[0_20px_80px_rgba(20,32,51,0.08)] backdrop-blur">
        <div className="grid gap-8 border-b border-border px-6 py-8 md:grid-cols-[1.3fr_0.7fr] md:px-10">
          <div className="space-y-6">
            <Badge variant="secondary" className="w-fit rounded-full px-4 py-1 text-[11px] uppercase tracking-[0.24em]">
              Next.js 16 + shadcn/ui
            </Badge>
            <div className="space-y-4">
              <h1 className="max-w-3xl text-4xl font-semibold tracking-tight md:text-6xl">
                Call center operations dashboard built around the Go API.
              </h1>
              <p className="max-w-2xl text-sm leading-7 text-muted-foreground md:text-base">
                The frontend is organized by feature slices, typed API access, and reusable primitives so it can
                expand into agent, supervisor, QA, and reporting workspaces without collapsing into one giant app.
              </p>
            </div>
            <div className="grid gap-3 sm:grid-cols-3">
              <Metric label="Frontend" value="Next.js 16" />
              <Metric label="API layer" value="Go + Chi" />
              <Metric label="Async core" value="Redis + RabbitMQ" />
            </div>
          </div>
          <div className="rounded-[1.6rem] border border-border bg-[linear-gradient(180deg,#18324f,#101b2d)] p-5 text-white">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-[0.24em] text-slate-300">System split</p>
                <p className="mt-2 text-2xl font-semibold">Frontend, API, Workers</p>
              </div>
              <div className="rounded-2xl bg-white/12 p-3 text-amber-300">
                <PhoneCall className="h-6 w-6" />
              </div>
            </div>
            <Separator className="my-5 bg-white/12" />
            <div className="grid gap-4">
              {pillars.map((pillar) => (
                <div key={pillar.title} className="rounded-2xl border border-white/10 bg-white/6 p-4">
                  <div className="flex items-center gap-3">
                    <pillar.icon className="h-4 w-4 text-amber-300" />
                    <h2 className="font-semibold">{pillar.title}</h2>
                  </div>
                  <p className="mt-2 text-sm leading-6 text-slate-300">{pillar.text}</p>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="px-6 py-8 md:px-10">{children}</div>
      </section>
    </main>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl border border-border bg-white/80 px-4 py-3">
      <p className="text-[11px] uppercase tracking-[0.2em] text-muted-foreground">{label}</p>
      <p className="mt-2 flex items-center gap-2 text-sm font-semibold text-foreground">
        <Activity className="h-4 w-4 text-primary" />
        {value}
      </p>
    </div>
  );
}
