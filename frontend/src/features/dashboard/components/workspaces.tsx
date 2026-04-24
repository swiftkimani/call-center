"use client";

import Link from "next/link";
import { useState } from "react";
import {
  BellRing,
  CheckCircle2,
  ChevronRight,
  Mic,
  MicOff,
  Phone,
  PhoneCall,
  PhoneForwarded,
  Users,
  Volume2,
  Waves,
} from "lucide-react";

import type { DashboardData, HealthItem } from "@/features/dashboard/types";
import { Badge } from "@/components/ui/badge";
import { Button, buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

const dispositions = ["Resolved", "Follow up", "Escalated", "No answer", "Converted"];

export function AgentWorkspacePage({ data }: { data: DashboardData }) {
  const [selectedDisposition, setSelectedDisposition] = useState(dispositions[0]);
  const [softphoneAction, setSoftphoneAction] = useState("Answer");

  return (
    <div className="space-y-4">
      <section className="grid gap-4 xl:grid-cols-[1.4fr_0.8fr]">
        <Card>
          <CardHeader>
            <div className="flex items-start justify-between gap-4">
              <div>
                <CardTitle className="text-xl">Live call</CardTitle>
                <CardDescription>Customer and call controls.</CardDescription>
              </div>
              <Badge variant="secondary" className="rounded-full px-3 py-1">
                {data.currentCall.queue}
              </Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="rounded-[1.5rem] bg-[linear-gradient(180deg,#20384a,#162734)] p-5 text-white">
              <p className="text-xs uppercase tracking-[0.16em] text-slate-300">Customer</p>
              <h2 className="mt-2 text-3xl font-semibold">{data.currentCall.callerName}</h2>
              <p className="mt-1 text-slate-300">{data.currentCall.callerNumber}</p>
              <div className="mt-4 flex flex-wrap gap-2">
                {data.currentCall.tags.slice(0, 3).map((tag) => (
                  <Badge key={tag} variant="secondary" className="rounded-full bg-white/10 text-white">
                    {tag}
                  </Badge>
                ))}
              </div>
              <p className="mt-4 text-sm leading-6 text-slate-200">{data.currentCall.summary}</p>
            </div>

            <div className="grid gap-3 sm:grid-cols-3">
              <MiniStat label="Call time" value="04:36" />
              <MiniStat label="Wait time" value={`${data.currentCall.waitSeconds}s`} />
              <MiniStat label="Disposition" value={selectedDisposition} />
            </div>

            <div>
              <p className="text-xs uppercase tracking-[0.16em] text-muted-foreground">Softphone</p>
              <div className="mt-3 grid grid-cols-3 gap-3 sm:grid-cols-6">
                <SoftphoneButton icon={Phone} label="Answer" active={softphoneAction === "Answer"} onClick={setSoftphoneAction} />
                <SoftphoneButton icon={MicOff} label="Mute" active={softphoneAction === "Mute"} onClick={setSoftphoneAction} />
                <SoftphoneButton icon={Volume2} label="Hold" active={softphoneAction === "Hold"} onClick={setSoftphoneAction} />
                <SoftphoneButton icon={PhoneForwarded} label="Transfer" active={softphoneAction === "Transfer"} onClick={setSoftphoneAction} />
                <SoftphoneButton icon={Mic} label="Whisper" active={softphoneAction === "Whisper"} onClick={setSoftphoneAction} />
                <SoftphoneButton icon={BellRing} label="DTMF" active={softphoneAction === "DTMF"} onClick={setSoftphoneAction} />
              </div>
            </div>

            <div>
              <p className="text-xs uppercase tracking-[0.16em] text-muted-foreground">Disposition</p>
              <div className="mt-3 flex flex-wrap gap-2">
                {dispositions.map((item) => (
                  <button
                    key={item}
                    type="button"
                    onClick={() => setSelectedDisposition(item)}
                    className={cn(
                      "rounded-full border px-3 py-2 text-xs font-semibold transition-colors",
                      selectedDisposition === item ? "border-primary/10 bg-primary text-primary-foreground" : "border-border bg-white hover:bg-secondary/60",
                    )}
                  >
                    {item}
                  </button>
                ))}
              </div>
            </div>

            <Button className="w-full sm:w-auto" onClick={() => setSoftphoneAction("Call ended")}>
              End call
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Call notes</CardTitle>
            <CardDescription>Only the context needed to finish the conversation.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <InfoRow label="Issue" value="Invoice appears duplicated before a business-line upgrade." />
            <InfoRow label="Next step" value="Confirm the fee, explain the adjustment, and offer a billing breakdown by email." />
            <InfoRow label="Last agent" value={data.customer.lastAgent} />
            <InfoRow label="Open ticket" value={data.customer.openTicket} />
          </CardContent>
        </Card>
      </section>
    </div>
  );
}

export function SupervisorCockpitPage({ data }: { data: DashboardData }) {
  const [liveAction, setLiveAction] = useState("No intervention selected");

  return (
    <div className="space-y-4">
      <PageHeader title="Supervisor" description={liveAction} />

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {data.queueSnapshots.map((queue) => (
          <Card key={queue.id}>
            <CardHeader className="pb-3">
              <CardDescription>{queue.name}</CardDescription>
              <CardTitle className="mt-1 text-3xl">{queue.waiting}</CardTitle>
            </CardHeader>
            <CardContent className="space-y-1 text-sm text-muted-foreground">
              <p>Longest wait {queue.longestWait}</p>
              <p>{queue.availableAgents} agents ready</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Agent actions</CardTitle>
          <CardDescription>Listen, whisper, or barge into a live call.</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
          {data.agentSnapshots.map((agent) => (
            <div key={agent.id} className="rounded-[1.1rem] border border-border bg-white/72 p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="font-semibold">{agent.name}</p>
                  <p className="text-sm text-muted-foreground">
                    {agent.queue} | Ext {agent.extension}
                  </p>
                </div>
                <Badge variant={agent.status === "busy" ? "default" : "secondary"}>{agent.status}</Badge>
              </div>
              <p className="mt-3 text-sm text-muted-foreground">{agent.liveCall}</p>
              <div className="mt-4 grid grid-cols-3 gap-2">
                <MiniAction icon={Volume2} label="Listen" onClick={() => setLiveAction(`Listening to ${agent.name}`)} />
                <MiniAction icon={Mic} label="Whisper" onClick={() => setLiveAction(`Whisper mode for ${agent.name}`)} />
                <MiniAction icon={PhoneCall} label="Barge" onClick={() => setLiveAction(`Barge request sent to ${agent.name}`)} />
              </div>
            </div>
          ))}
        </CardContent>
      </Card>
    </div>
  );
}

export function CampaignManagerPage({ data }: { data: DashboardData }) {
  const [selectedCampaignId, setSelectedCampaignId] = useState(data.campaigns[0]?.id ?? null);
  const selectedCampaign = data.campaigns.find((campaign) => campaign.id === selectedCampaignId) ?? data.campaigns[0] ?? null;

  return (
    <div className="space-y-4">
      <PageHeader title="Campaigns" description="Select a campaign to inspect or open." />

      <section className="grid gap-4 xl:grid-cols-[0.95fr_1.05fr]">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Campaign list</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {data.campaigns.map((campaign) => (
              <button
                key={campaign.id}
                type="button"
                onClick={() => setSelectedCampaignId(campaign.id)}
                className={cn(
                  "flex w-full items-center justify-between rounded-[1rem] border px-4 py-4 text-left transition-colors",
                  selectedCampaignId === campaign.id ? "border-primary/15 bg-primary/[0.05]" : "border-border bg-white/70 hover:bg-white",
                )}
              >
                <div>
                  <p className="font-semibold">{campaign.name}</p>
                  <p className="text-sm text-muted-foreground">{campaign.scheduled_at ?? "Not scheduled"}</p>
                </div>
                <Badge variant="secondary">{campaign.status ?? "draft"}</Badge>
              </button>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Selected campaign</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {selectedCampaign && (
              <>
                <InfoRow label="Name" value={selectedCampaign.name} />
                <InfoRow label="Status" value={selectedCampaign.status ?? "draft"} />
                <InfoRow label="Scheduled" value={selectedCampaign.scheduled_at ?? "Not scheduled"} />
              </>
            )}
            <Button>Create campaign</Button>
          </CardContent>
        </Card>
      </section>
    </div>
  );
}

export function ReportsPage({ data }: { data: DashboardData }) {
  return (
    <div className="space-y-4">
      <PageHeader title="Reports" description="Core metrics and export shortcuts." />

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {data.reportCards.map((card) => (
          <Card key={card.label}>
            <CardHeader className="pb-3">
              <CardDescription>{card.label}</CardDescription>
              <CardTitle className="mt-1 text-3xl">{card.value}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-emerald-700">{card.delta}</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Exports</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-2">
          <AssistCard icon={Users} title="Agent performance" text="Per-agent snapshots." href="/supervisor" />
          <AssistCard icon={Waves} title="Queue trend pack" text="Wait times and abandonment." href="/reports" />
          <AssistCard icon={CheckCircle2} title="Campaign outcomes" text="Conversion and retries." href="/campaigns" />
        </CardContent>
      </Card>
    </div>
  );
}

export function SettingsPage({ data }: { data: DashboardData }) {
  const [selectedQueueId, setSelectedQueueId] = useState(data.queueSnapshots[0]?.id ?? null);
  const selectedQueue = data.queueSnapshots.find((queue) => queue.id === selectedQueueId) ?? data.queueSnapshots[0] ?? null;

  return (
    <div className="space-y-4">
      <PageHeader title="Settings" description="Queue configuration only." />

      <section className="grid gap-4 xl:grid-cols-[0.9fr_1.1fr]">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Queues</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {data.queueSnapshots.map((queue) => (
              <button
                key={queue.id}
                type="button"
                onClick={() => setSelectedQueueId(queue.id)}
                className={cn(
                  "flex w-full items-center justify-between rounded-[1rem] border px-4 py-4 text-left transition-colors",
                  selectedQueueId === queue.id ? "border-primary/15 bg-primary/[0.05]" : "border-border bg-white/70 hover:bg-white",
                )}
              >
                <div>
                  <p className="font-semibold">{queue.name}</p>
                  <p className="text-sm text-muted-foreground">{queue.availableAgents} agents ready</p>
                </div>
                <ChevronRight className="h-4 w-4 text-muted-foreground" />
              </button>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Queue detail</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {selectedQueue && (
              <>
                <InfoRow label="Queue" value={selectedQueue.name} />
                <InfoRow label="Waiting callers" value={`${selectedQueue.waiting}`} />
                <InfoRow label="Available agents" value={`${selectedQueue.availableAgents}`} />
                <InfoRow label="Longest wait" value={selectedQueue.longestWait} />
              </>
            )}
            {data.platformHealth.slice(0, 2).map((item) => (
              <HealthRow key={item.name} item={item} />
            ))}
          </CardContent>
        </Card>
      </section>
    </div>
  );
}

function PageHeader({ title, description }: { title: string; description: string }) {
  return (
    <div className="rounded-[1rem] border border-border bg-white/68 px-4 py-4">
      <h2 className="text-xl font-semibold tracking-tight">{title}</h2>
      <p className="mt-1 text-sm text-muted-foreground">{description}</p>
    </div>
  );
}

function MiniStat({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-[1rem] border border-border bg-white/70 px-4 py-3">
      <p className="text-xs uppercase tracking-[0.16em] text-muted-foreground">{label}</p>
      <p className="mt-2 text-lg font-semibold">{value}</p>
    </div>
  );
}

function SoftphoneButton({
  icon: Icon,
  label,
  active,
  onClick,
}: {
  icon: typeof Phone;
  label: string;
  active: boolean;
  onClick: (label: string) => void;
}) {
  return (
    <button
      type="button"
      onClick={() => onClick(label)}
      className={cn(
        "rounded-[1rem] border px-3 py-4 text-center text-sm font-semibold transition-colors",
        active ? "border-primary/10 bg-primary text-primary-foreground" : "border-border bg-white hover:bg-secondary/60",
      )}
    >
      <Icon className="mx-auto h-5 w-5" />
      <span className="mt-2 block">{label}</span>
    </button>
  );
}

function MiniAction({ icon: Icon, label, onClick }: { icon: typeof Phone; label: string; onClick: () => void }) {
  return (
    <button type="button" onClick={onClick} className="rounded-xl border border-border bg-white px-2 py-2 text-xs font-semibold transition-colors hover:bg-secondary/60">
      <Icon className="mx-auto h-4 w-4 text-primary" />
      <span className="mt-1 block">{label}</span>
    </button>
  );
}

function AssistCard({
  icon: Icon,
  title,
  text,
  href,
}: {
  icon: typeof Users;
  title: string;
  text: string;
  href: string;
}) {
  return (
    <Link href={href} className="block rounded-[1rem] border border-border bg-white/72 p-4 transition-colors hover:bg-white">
      <Icon className="h-5 w-5 text-primary" />
      <p className="mt-3 font-semibold">{title}</p>
      <p className="mt-1 text-sm text-muted-foreground">{text}</p>
      <span className={cn(buttonVariants({ variant: "outline", size: "sm" }), "mt-4 inline-flex")}>Open</span>
    </Link>
  );
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-[1rem] border border-border bg-white/72 px-4 py-4">
      <p className="text-xs uppercase tracking-[0.16em] text-muted-foreground">{label}</p>
      <p className="mt-2 text-sm leading-6 text-foreground">{value}</p>
    </div>
  );
}

function HealthRow({ item }: { item: HealthItem }) {
  return (
    <div className="rounded-[1rem] border border-border bg-white/72 p-4">
      <div className="flex items-center justify-between gap-3">
        <p className="font-semibold">{item.name}</p>
        <Badge variant={item.status === "healthy" ? "secondary" : "default"}>{item.status}</Badge>
      </div>
      <p className="mt-2 text-sm text-muted-foreground">{item.detail}</p>
    </div>
  );
}
