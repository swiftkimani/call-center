import { AppShell } from "@/components/layout/app-shell";
import { DashboardPage } from "@/features/dashboard/components/dashboard-page";
import { getDashboardData } from "@/features/dashboard/server/get-dashboard-data";

export default async function HomePage() {
  const data = await getDashboardData();

  return (
    <AppShell>
      <DashboardPage data={data} />
    </AppShell>
  );
}
