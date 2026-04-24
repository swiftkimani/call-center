import "server-only";

import { backendGateway } from "@/lib/http/backend-gateway";
import type { Campaign, CallRecord } from "@/lib/http/types";
import type { DashboardData } from "@/features/dashboard/types";

export async function getDashboardData(): Promise<DashboardData> {
  const [health, me, campaigns, calls] = await Promise.all([
    backendGateway.health(),
    backendGateway.me.get().catch(() => null),
    backendGateway.campaigns.list({ limit: 5 }).catch(() => null),
    backendGateway.calls.list({ limit: 6 }).catch(() => null),
  ]);

  const resolvedCampaigns = campaigns?.success && Array.isArray(campaigns.data) ? campaigns.data : fallbackCampaigns;
  const resolvedCalls = calls?.success && Array.isArray(calls.data) ? calls.data : fallbackCalls;
  const resolvedMe = me?.success ? me.data : null;

  return {
    health,
    me: resolvedMe,
    campaigns: resolvedCampaigns,
    calls: resolvedCalls,
    currentCall: {
      id: resolvedCalls[0]?.id ?? "call-live-001",
      callerName: "Jane Mwangi",
      callerNumber: "+254 712 345 678",
      queue: "Sales",
      waitSeconds: 12,
      sentiment: "Calm but urgent",
      tags: ["VIP", "Repeat caller", "Nairobi"],
      summary: "Wants to upgrade a business line and confirm why the last invoice posted twice.",
    },
    customer: {
      id: resolvedCalls[0]?.customer_id ?? "cust-001",
      name: "Jane Mwangi",
      plan: "Business Premium",
      city: "Nairobi",
      language: "English / Swahili",
      lastAgent: "Brian Ouma",
      openTicket: "Billing reconciliation pending",
      balanceNote: "KES 8,400 invoice disputed on the last cycle.",
      lastOrders: ["Line activation - Apr 14", "SIM replacement - Apr 05", "Roaming add-on - Mar 27"],
    },
    activity: [
      {
        id: "act-1",
        title: "Incoming VIP call routed",
        detail: "Sales queue sent the customer to the longest-idle qualified agent.",
        time: "Now",
        tone: "good",
      },
      {
        id: "act-2",
        title: "SLA risk on Retention queue",
        detail: "2 callers have waited longer than 90 seconds.",
        time: "2 min ago",
        tone: "warn",
      },
      {
        id: "act-3",
        title: "Campaign contacts imported",
        detail: `${resolvedCampaigns[0]?.name ?? "Winback April"} received 2,400 new contacts.`,
        time: "18 min ago",
        tone: "neutral",
      },
      {
        id: "act-4",
        title: "Recording worker healthy",
        detail: "MinIO uploads are completing within normal latency.",
        time: "23 min ago",
        tone: "good",
      },
    ],
    queueSnapshots: [
      { id: "q-sales", name: "Sales", waiting: 7, slaBreaches: 1, longestWait: "01:28", availableAgents: 6 },
      { id: "q-support", name: "Support", waiting: 3, slaBreaches: 0, longestWait: "00:42", availableAgents: 11 },
      { id: "q-retention", name: "Retention", waiting: 5, slaBreaches: 2, longestWait: "02:11", availableAgents: 4 },
      { id: "q-priority", name: "Priority", waiting: 1, slaBreaches: 0, longestWait: "00:09", availableAgents: 3 },
    ],
    agentSnapshots: [
      { id: "a-1", name: "Amina Hassan", extension: "201", status: "available", queue: "Priority", occupancy: "61%", liveCall: "Idle" },
      { id: "a-2", name: "Brian Ouma", extension: "202", status: "busy", queue: "Sales", occupancy: "87%", liveCall: "Jane Mwangi" },
      { id: "a-3", name: "Cynthia Wairimu", extension: "203", status: "wrap_up", queue: "Support", occupancy: "72%", liveCall: "Post-call notes" },
      { id: "a-4", name: "David Kiptoo", extension: "204", status: "break", queue: "Retention", occupancy: "33%", liveCall: "On break" },
      { id: "a-5", name: "Esther Njeri", extension: "205", status: "available", queue: "Sales", occupancy: "58%", liveCall: "Idle" },
      { id: "a-6", name: "Faith Atieno", extension: "206", status: "busy", queue: "Support", occupancy: "79%", liveCall: "Account update" },
    ],
    reportCards: [
      { label: "Answer rate", value: "94.2%", delta: "+2.1% vs yesterday" },
      { label: "Avg handle time", value: "04:36", delta: "-00:18 vs yesterday" },
      { label: "Agent utilization", value: "71%", delta: "+5 pts vs last week" },
      { label: "Campaign conversion", value: "18.4%", delta: "+3.2 pts vs previous run" },
    ],
    platformHealth: [
      { name: "API", status: health.status, detail: health.detail },
      { name: "Redis realtime", status: "healthy", detail: "Presence fan-out within target latency." },
      { name: "RabbitMQ", status: "healthy", detail: "Outbound and recording queues are draining normally." },
      { name: "MinIO", status: "healthy", detail: "Recording uploads and signed playback URLs are available." },
    ],
  };
}

const fallbackCampaigns: Campaign[] = [
  { id: "cmp-1", name: "April Winback", status: "running", scheduled_at: "2026-04-24T09:00:00Z" },
  { id: "cmp-2", name: "SME Upgrade Push", status: "scheduled", scheduled_at: "2026-04-24T13:30:00Z" },
  { id: "cmp-3", name: "Collections Morning Block", status: "draft", scheduled_at: "2026-04-25T08:00:00Z" },
];

const fallbackCalls: CallRecord[] = [
  { id: "call-1", agent_id: "a-2", customer_id: "cust-001", direction: "inbound", status: "in_progress" },
  { id: "call-2", agent_id: "a-6", customer_id: "cust-002", direction: "outbound", status: "ringing" },
  { id: "call-3", agent_id: "a-3", customer_id: "cust-003", direction: "inbound", status: "completed" },
  { id: "call-4", agent_id: "a-5", customer_id: "cust-004", direction: "outbound", status: "queued" },
];
