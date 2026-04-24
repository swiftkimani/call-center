import "server-only";

import { getApiBaseUrl, getApiToken } from "@/lib/env";
import type { HealthCheck } from "@/lib/http/types";

type RequestInitWithNext = RequestInit & {
  next?: {
    revalidate?: number;
    tags?: string[];
  };
};

async function request<T>(path: string, init?: RequestInitWithNext): Promise<T> {
  const token = getApiToken();
  const headers = new Headers(init?.headers);

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${getApiBaseUrl()}${path}`, {
    ...init,
    headers,
  });

  if (!response.ok) {
    throw new Error(`API request failed for ${path}: ${response.status}`);
  }

  const contentType = response.headers.get("content-type") ?? "";
  if (contentType.includes("application/json")) {
    return (await response.json()) as T;
  }

  return (await response.text()) as T;
}

export const apiClient = {
  health: async (): Promise<HealthCheck> => {
    try {
      const result = await request<string>("/healthz", {
        cache: "no-store",
      });

      return {
        status: result === "ok" ? "healthy" : "unknown",
        detail: `GET /healthz -> ${result}`,
      };
    } catch {
      return {
        status: "offline",
        detail: `Could not reach ${getApiBaseUrl()}/healthz`,
      };
    }
  },
};

export async function optionalApiClient<T>(path: string, init?: RequestInitWithNext): Promise<T | null> {
  try {
    return await request<T>(path, init);
  } catch {
    return null;
  }
}
