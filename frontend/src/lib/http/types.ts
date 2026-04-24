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
};
