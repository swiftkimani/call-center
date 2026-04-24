import { AgentWorkspacePage } from "@/features/dashboard/components/workspaces";
import { getDashboardData } from "@/features/dashboard/server/get-dashboard-data";

export default async function AgentPage() {
  const data = await getDashboardData();

  return <AgentWorkspacePage data={data} />;
}
