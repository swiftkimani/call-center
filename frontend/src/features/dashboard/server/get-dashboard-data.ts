import "server-only";

import { apiClient, optionalApiClient } from "@/lib/http/api-client";
import type { ApiResponse, Campaign, CallRecord, Me } from "@/lib/http/types";
import type { DashboardData } from "@/features/dashboard/types";

export async function getDashboardData(): Promise<DashboardData> {
  const [health, me, campaigns, calls] = await Promise.all([
    apiClient.health(),
    optionalApiClient<ApiResponse<Me>>("/api/v1/me"),
    optionalApiClient<ApiResponse<Campaign[]>>("/api/v1/campaigns?limit=5"),
    optionalApiClient<ApiResponse<CallRecord[]>>("/api/v1/calls?limit=6"),
  ]);

  return {
    health,
    me: me?.success ? me.data : null,
    campaigns: campaigns?.success && Array.isArray(campaigns.data) ? campaigns.data : [],
    calls: calls?.success && Array.isArray(calls.data) ? calls.data : [],
  };
}
