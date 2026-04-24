const DEFAULT_API_BASE_URL = "http://localhost:8080";

export function getApiBaseUrl() {
  return process.env.CALLCENTER_API_BASE_URL ?? process.env.NEXT_PUBLIC_API_BASE_URL ?? DEFAULT_API_BASE_URL;
}

export function getApiToken() {
  return process.env.CALLCENTER_API_TOKEN;
}
