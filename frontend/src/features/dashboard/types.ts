import type { Campaign, CallRecord, HealthCheck, Me } from "@/lib/http/types";

export type DashboardData = {
  health: HealthCheck;
  me: Me | null;
  campaigns: Campaign[];
  calls: CallRecord[];
};
