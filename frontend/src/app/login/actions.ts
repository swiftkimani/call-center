"use server";

import { redirect } from "next/navigation";

import { backendGateway } from "@/lib/http/backend-gateway";
import { setSession } from "@/lib/auth/session";
import { getApiBaseUrl } from "@/lib/env";

export type LoginActionState = {
  error?: string;
};

export async function loginAction(_: LoginActionState, formData: FormData): Promise<LoginActionState> {
  const email = String(formData.get("email") ?? "").trim();
  const password = String(formData.get("password") ?? "");

  if (!email || !password) {
    return { error: "Email and password are required." };
  }

  try {
    const response = await backendGateway.auth.login({ email, password });

    if (!response.success || !response.data) {
      return { error: response.error ?? "Login failed." };
    }

    await setSession({
      accessToken: response.data.access_token,
      refreshToken: response.data.refresh_token,
      role: response.data.role,
    });
  } catch {
    return { error: `Could not reach the backend at ${getApiBaseUrl()}.` };
  }

  redirect("/agent");
}
