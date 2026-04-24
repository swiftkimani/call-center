import "server-only";

import { cookies } from "next/headers";

const ACCESS_TOKEN_COOKIE = "callcenter_access_token";
const REFRESH_TOKEN_COOKIE = "callcenter_refresh_token";
const ROLE_COOKIE = "callcenter_role";

const cookieOptions = {
  httpOnly: true,
  sameSite: "lax" as const,
  secure: process.env.NODE_ENV === "production",
  path: "/",
  maxAge: 60 * 60 * 24 * 7,
};

export async function getSessionAccessToken() {
  const store = await cookies();
  return store.get(ACCESS_TOKEN_COOKIE)?.value ?? process.env.CALLCENTER_API_TOKEN ?? null;
}

export async function getSessionRole() {
  const store = await cookies();
  return store.get(ROLE_COOKIE)?.value ?? null;
}

export async function hasSession() {
  return Boolean(await getSessionAccessToken());
}

export async function setSession(tokens: { accessToken: string; refreshToken: string; role: string }) {
  const store = await cookies();
  store.set(ACCESS_TOKEN_COOKIE, tokens.accessToken, cookieOptions);
  store.set(REFRESH_TOKEN_COOKIE, tokens.refreshToken, cookieOptions);
  store.set(ROLE_COOKIE, tokens.role, { ...cookieOptions, httpOnly: false });
}

export async function clearSession() {
  const store = await cookies();
  store.delete(ACCESS_TOKEN_COOKIE);
  store.delete(REFRESH_TOKEN_COOKIE);
  store.delete(ROLE_COOKIE);
}
