import type { ComponentType } from "react";
import { Activity, CheckCircle2, Clock3, LayoutDashboard, ShieldAlert, Waves } from "lucide-react";

import type { DashboardData } from "@/features/dashboard/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

const kpis = (data: DashboardData) => [
  {
    label: "API reachability",
    value: data.health.status,
    detail: data.health.detail,
    icon: CheckCircle2,
  },
  {
    label: "Campaigns loaded",
    value: String(data.campaigns.length),
    detail: "Server-rendered from /api/v1/campaigns",
    icon: LayoutDashboard,
  },
  {
    label: "Calls loaded",
    value: String(data.calls.length),
    detail: "Server-rendered from /api/v1/calls",
    icon: Waves,
  },
  {
    label: "Current identity",
    value: data.me?.full_name ?? "Token not configured",
    detail: data.me ? data.me.role : "Set CALLCENTER_API_TOKEN to unlock protected widgets",
    icon: Activity,
  },
];

export function DashboardPage({ data }: { data: DashboardData }) {
  return (
    <div className="space-y-8">
      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {kpis(data).map((item) => (
          <Card key={item.label} className="border-border bg-white">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <div>
                <CardDescription>{item.label}</CardDescription>
                <CardTitle className="mt-2 text-2xl">{item.value}</CardTitle>
              </div>
              <div className="rounded-2xl bg-secondary p-3 text-primary">
                <item.icon className="h-5 w-5" />
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm leading-6 text-muted-foreground">{item.detail}</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <section className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
        <Card className="border-border bg-white">
          <CardHeader>
            <div className="flex items-center justify-between gap-4">
              <div>
                <CardTitle>Campaign pipeline</CardTitle>
                <CardDescription>Server-rendered campaign snapshot with graceful fallback when protected endpoints reject requests.</CardDescription>
              </div>
              <Badge variant="outline">{data.campaigns.length} visible</Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-3">
            {data.campaigns.length > 0 ? (
              data.campaigns.map((campaign) => (
                <div key={campaign.id} className="rounded-2xl border border-border bg-secondary/30 p-4">
                  <div className="flex flex-wrap items-center justify-between gap-3">
                    <div>
                      <p className="font-semibold">{campaign.name}</p>
                      <p className="text-sm text-muted-foreground">
                        Status: {campaign.status ?? "pending"} • Scheduled: {campaign.scheduled_at ?? "not scheduled"}
                      </p>
                    </div>
                    <Badge variant="secondary">{campaign.status ?? "draft"}</Badge>
                  </div>
                </div>
              ))
            ) : (
              <EmptyState
                icon={Clock3}
                title="No campaigns available from the current API context."
                body="If your backend requires auth, add CALLCENTER_API_TOKEN to frontend/.env.local and reload."
              />
            )}
          </CardContent>
        </Card>

        <Card className="border-border bg-[linear-gradient(180deg,#16314d,#102237)] text-white">
          <CardHeader>
            <CardTitle>Scale blueprint</CardTitle>
            <CardDescription className="text-slate-300">
              Frontend and backend are split for vertical ownership and horizontal runtime scaling.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm leading-6 text-slate-200">
            <div className="rounded-2xl bg-white/8 p-4">
              <p className="font-semibold text-white">Frontend</p>
              <p>App Router, feature slices, server-first data fetching, and shared UI primitives.</p>
            </div>
            <div className="rounded-2xl bg-white/8 p-4">
              <p className="font-semibold text-white">API</p>
              <p>Stateless Go HTTP service behind a load balancer, backed by Redis, RabbitMQ, and Postgres.</p>
            </div>
            <div className="rounded-2xl bg-white/8 p-4">
              <p className="font-semibold text-white">Workers</p>
              <p>Separate deployment units for campaigns, recordings, notifications, and future ML enrichments.</p>
            </div>
            <Button className="w-full bg-white text-slate-900 hover:bg-slate-100">
              <span>See `docs/architecture/scalable-structure.md`</span>
            </Button>
          </CardContent>
        </Card>
      </section>

      <section className="grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
        <Card className="border-border bg-white">
          <CardHeader>
            <CardTitle>Authenticated context</CardTitle>
            <CardDescription>Protected API data is optional but ready for agent and supervisor dashboards.</CardDescription>
          </CardHeader>
          <CardContent>
            {data.me ? (
              <div className="space-y-3">
                <div className="rounded-2xl bg-secondary/50 p-4">
                  <p className="text-xs uppercase tracking-[0.2em] text-muted-foreground">User</p>
                  <p className="mt-2 text-lg font-semibold">{data.me.full_name}</p>
                  <p className="text-sm text-muted-foreground">{data.me.email}</p>
                </div>
                <div className="grid gap-3 md:grid-cols-2">
                  <InfoBlock label="Role" value={data.me.role} />
                  <InfoBlock label="Agent status" value={data.me.agent?.status ?? "n/a"} />
                </div>
              </div>
            ) : (
              <EmptyState
                icon={ShieldAlert}
                title="Protected widgets are idle."
                body="The frontend is using the API client correctly, but it has no bearer token configured for /api/v1/me."
              />
            )}
          </CardContent>
        </Card>

        <Card className="border-border bg-white">
          <CardHeader>
            <CardTitle>Recent calls</CardTitle>
            <CardDescription>Server-rendered snapshot from the Go API. Replace with TanStack Query + WebSocket hydration for live views.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {data.calls.length > 0 ? (
              data.calls.map((call) => (
                <div key={call.id} className="flex flex-col gap-2 rounded-2xl border border-border bg-secondary/30 p-4 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="font-semibold">{call.direction ?? "unknown"} call</p>
                    <p className="text-sm text-muted-foreground">
                      Agent: {call.agent_id ?? "unassigned"} • Customer: {call.customer_id ?? "unknown"}
                    </p>
                  </div>
                  <Badge variant="outline">{call.status ?? "pending"}</Badge>
                </div>
              ))
            ) : (
              <EmptyState
                icon={Clock3}
                title="No calls returned."
                body="The UI is ready; populate sample data or seed the backend to see live call rows here."
              />
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  );
}

function EmptyState({
  icon: Icon,
  title,
  body,
}: {
  icon: ComponentType<{ className?: string }>;
  title: string;
  body: string;
}) {
  return (
    <div className="rounded-[1.5rem] border border-dashed border-border bg-secondary/30 p-6">
      <Icon className="h-5 w-5 text-primary" />
      <p className="mt-4 font-semibold">{title}</p>
      <p className="mt-2 text-sm leading-6 text-muted-foreground">{body}</p>
    </div>
  );
}

function InfoBlock({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl border border-border bg-secondary/20 p-4">
      <p className="text-xs uppercase tracking-[0.18em] text-muted-foreground">{label}</p>
      <p className="mt-2 font-semibold">{value}</p>
    </div>
  );
}
