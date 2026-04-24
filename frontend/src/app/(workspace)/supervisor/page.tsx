import { SupervisorCockpitPage } from "@/features/dashboard/components/workspaces";
import { getDashboardData } from "@/features/dashboard/server/get-dashboard-data";

export default async function SupervisorPage() {
  const data = await getDashboardData();

  return <SupervisorCockpitPage data={data} />;
}
