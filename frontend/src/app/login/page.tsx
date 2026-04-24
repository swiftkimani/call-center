import { redirect } from "next/navigation";

import { LoginScreen } from "@/components/auth/login-screen";
import { hasSession } from "@/lib/auth/session";

export default async function LoginPage() {
  redirect((await hasSession()) ? "/agent" : "/");
}
