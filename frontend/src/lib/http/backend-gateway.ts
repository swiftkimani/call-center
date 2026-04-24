import "server-only";

import { apiClient } from "@/lib/http/api-client";
import type {
  ApiResponse,
  CallDispositionRequest,
  CallRecord,
  Campaign,
  CampaignImportResponse,
  CampaignContactsImportRequest,
  CreateCampaignRequest,
  CustomerRecord,
  DailyReportSummary,
  LoginRequest,
  LoginResponse,
  Me,
  OutboundCallRequest,
  QueueLiveSnapshot,
  RecordingResponse,
  RefreshRequest,
  SetAgentStatusRequest,
  StatusResponse,
  UpdateCampaignStatusRequest,
  UpdateCustomerRequest,
} from "@/lib/http/types";

export const backendGateway = {
  health: apiClient.health,

  auth: {
    login: (payload: LoginRequest) => apiClient.post<ApiResponse<LoginResponse>>("/api/v1/auth/login", payload),
    refresh: (payload: RefreshRequest) => apiClient.post<ApiResponse<LoginResponse>>("/api/v1/auth/refresh", payload),
    logout: () => apiClient.post<ApiResponse<null>>("/api/v1/auth/logout"),
  },

  me: {
    get: () => apiClient.get<ApiResponse<Me>>("/api/v1/me"),
  },

  agents: {
    setStatus: (agentId: string, payload: SetAgentStatusRequest) => apiClient.post<ApiResponse<StatusResponse>>(`/api/v1/agents/${agentId}/status`, payload),
    heartbeat: () => apiClient.post<void>("/api/v1/agents/heartbeat"),
  },

  calls: {
    list: (params?: { limit?: number; offset?: number }) => apiClient.get<ApiResponse<CallRecord[]>>(withSearch("/api/v1/calls", params)),
    initiateOutbound: (payload: OutboundCallRequest) => apiClient.post<ApiResponse<CallRecord>>("/api/v1/calls/outbound", payload),
    saveDisposition: (callId: string, payload: CallDispositionRequest) => apiClient.post<ApiResponse<null>>(`/api/v1/calls/${callId}/disposition`, payload),
    getRecording: (callId: string) => apiClient.get<ApiResponse<RecordingResponse>>(`/api/v1/calls/${callId}/recording`),
  },

  campaigns: {
    list: (params?: { limit?: number; offset?: number }) => apiClient.get<ApiResponse<Campaign[]>>(withSearch("/api/v1/campaigns", params)),
    get: (campaignId: string) => apiClient.get<ApiResponse<Campaign>>(`/api/v1/campaigns/${campaignId}`),
    create: (payload: CreateCampaignRequest) => apiClient.post<ApiResponse<Campaign>>("/api/v1/campaigns", payload),
    importContacts: (campaignId: string, payload: CampaignContactsImportRequest) =>
      apiClient.post<ApiResponse<CampaignImportResponse>>(`/api/v1/campaigns/${campaignId}/contacts`, payload),
    updateStatus: (campaignId: string, payload: UpdateCampaignStatusRequest) =>
      apiClient.post<ApiResponse<StatusResponse>>(`/api/v1/campaigns/${campaignId}/status`, payload),
  },

  customers: {
    search: (params?: { q?: string; limit?: number; offset?: number }) => apiClient.get<ApiResponse<CustomerRecord[]>>(withSearch("/api/v1/customers", params)),
    get: (customerId: string) => apiClient.get<ApiResponse<CustomerRecord>>(`/api/v1/customers/${customerId}`),
    update: (customerId: string, payload: UpdateCustomerRequest) => apiClient.put<ApiResponse<CustomerRecord>>(`/api/v1/customers/${customerId}`, payload),
  },

  supervisor: {
    queueLive: (queueId: string) => apiClient.get<ApiResponse<QueueLiveSnapshot>>(`/api/v1/queues/${queueId}/live`),
    whisper: (callId: string) => apiClient.post<ApiResponse<null>>(`/api/v1/supervisor/${callId}/whisper`),
    barge: (callId: string) => apiClient.post<ApiResponse<StatusResponse>>(`/api/v1/supervisor/${callId}/barge`),
  },

  reports: {
    daily: (params?: { date?: string }) => apiClient.get<ApiResponse<DailyReportSummary>>(withSearch("/api/v1/reports/daily", params)),
    dailyCsvUrl: (date?: string) => withSearch("/api/v1/reports/daily", { date, format: "csv" }),
  },

  realtime: {
    agentSocketPath: "/ws/agent",
    supervisorSocketPath: "/ws/supervisor",
  },
};

function withSearch(path: string, params?: Record<string, string | number | undefined>) {
  if (!params) return path;

  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === "") continue;
    search.set(key, String(value));
  }

  const query = search.toString();
  return query ? `${path}?${query}` : path;
}
