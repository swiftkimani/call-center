import { SettingsPage } from "@/features/dashboard/components/workspaces";
import { getDashboardData } from "@/features/dashboard/server/get-dashboard-data";

export default async function SettingsRoutePage() {
  const data = await getDashboardData();

  return <SettingsPage data={data} />;
}
