import { ReportsPage } from "@/features/dashboard/components/workspaces";
import { getDashboardData } from "@/features/dashboard/server/get-dashboard-data";

export default async function ReportsRoutePage() {
  const data = await getDashboardData();

  return <ReportsPage data={data} />;
}
