export type ApiResponse<T> = {
  success: boolean;
  data: T;
  error?: string;
};

export type HealthCheck = {
  status: string;
  detail: string;
};

export type AgentProfile = {
  id: string;
  extension: string;
  skills: string[];
  status: string;
};

export type Me = {
  id: string;
  email: string;
  full_name: string;
  role: string;
  agent?: AgentProfile;
};

export type Campaign = {
  id: string;
  name: string;
  status?: string | null;
  scheduled_at?: string | null;
};

export type CallRecord = {
  id: string;
  agent_id?: string | null;
  customer_id?: string | null;
  direction?: string | null;
  status?: string | null;
  recording_url?: string | null;
};

export type CustomerRecord = {
  id: string;
  full_name: string;
  phone_e164?: string | null;
  email?: string | null;
  tags?: string[] | null;
};

export type QueueLiveSnapshot = {
  queue_id: string;
  name: string;
  description: string;
  skills_required: string[];
  max_wait_seconds: number;
  sla_seconds: number;
  waiting: number;
  oldest_wait_seconds: number;
  sla_breach_count: number;
};

export type RecordingResponse = {
  url: string;
};

export type StatusResponse = {
  status: string;
};

export type CampaignImportResponse = {
  campaign_id: string;
  imported: number;
};

export type DailyAgentSummary = {
  agent_id: string;
  total_calls: number;
  completed_calls: number;
  abandoned_calls: number;
  avg_talk_seconds: number;
  avg_wait_seconds: number;
  total_cost_cents: number;
};

export type DailyReportSummary = {
  date: string;
  agents: DailyAgentSummary[];
};

export type LoginResponse = {
  access_token: string;
  refresh_token: string;
  role: string;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type RefreshRequest = {
  refresh_token: string;
};

export type SetAgentStatusRequest = {
  status: string;
};

export type OutboundCallRequest = {
  customer_id: string;
  customer_phone: string;
};

export type CallDispositionRequest = {
  category: string;
  notes?: string;
};

export type CreateCampaignRequest = {
  name: string;
  scheduled_at?: string;
};

export type CampaignContactsImportRequest = {
  contacts: Array<{
    customer_id: string;
  }>;
};

export type UpdateCampaignStatusRequest = {
  status: string;
};

export type UpdateCustomerRequest = {
  full_name: string;
  email?: string;
  tags?: string[];
};
