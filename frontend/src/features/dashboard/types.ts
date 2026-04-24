import type { Campaign, CallRecord, HealthCheck, Me } from "@/lib/http/types";

export type WorkspaceCall = {
  id: string;
  callerName: string;
  callerNumber: string;
  queue: string;
  waitSeconds: number;
  sentiment: string;
  tags: string[];
  summary: string;
};

export type CustomerInsight = {
  id: string;
  name: string;
  plan: string;
  city: string;
  language: string;
  lastAgent: string;
  openTicket: string;
  balanceNote: string;
  lastOrders: string[];
};

export type ActivityItem = {
  id: string;
  title: string;
  detail: string;
  time: string;
  tone: "neutral" | "good" | "warn";
};

export type QueueSnapshot = {
  id: string;
  name: string;
  waiting: number;
  slaBreaches: number;
  longestWait: string;
  availableAgents: number;
};

export type AgentSnapshot = {
  id: string;
  name: string;
  extension: string;
  status: string;
  queue: string;
  occupancy: string;
  liveCall: string;
};

export type ReportCard = {
  label: string;
  value: string;
  delta: string;
};

export type HealthItem = {
  name: string;
  status: string;
  detail: string;
};

export type DashboardData = {
  health: HealthCheck;
  me: Me | null;
  campaigns: Campaign[];
  calls: CallRecord[];
  currentCall: WorkspaceCall;
  customer: CustomerInsight;
  activity: ActivityItem[];
  queueSnapshots: QueueSnapshot[];
  agentSnapshots: AgentSnapshot[];
  reportCards: ReportCard[];
  platformHealth: HealthItem[];
};
