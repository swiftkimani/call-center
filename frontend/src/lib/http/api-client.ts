import "server-only";

import { getSessionAccessToken } from "@/lib/auth/session";
import { getApiBaseUrl } from "@/lib/env";
import type { HealthCheck } from "@/lib/http/types";

type RequestInitWithNext = RequestInit & {
  next?: {
    revalidate?: number;
    tags?: string[];
  };
};

async function request<T>(path: string, init?: RequestInitWithNext): Promise<T> {
  const token = await getSessionAccessToken();
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
  get: <T>(path: string, init?: RequestInitWithNext) => request<T>(path, init),
  post: <T>(path: string, body?: unknown, init?: RequestInitWithNext) =>
    request<T>(path, {
      ...init,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...Object.fromEntries(new Headers(init?.headers).entries()),
      },
      body: body === undefined ? undefined : JSON.stringify(body),
    }),
  put: <T>(path: string, body?: unknown, init?: RequestInitWithNext) =>
    request<T>(path, {
      ...init,
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        ...Object.fromEntries(new Headers(init?.headers).entries()),
      },
      body: body === undefined ? undefined : JSON.stringify(body),
    }),
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
