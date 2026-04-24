import { redirect } from "next/navigation";

import { hasSession } from "@/lib/auth/session";

export default async function LoginPage() {
  redirect((await hasSession()) ? "/agent" : "/");
}
