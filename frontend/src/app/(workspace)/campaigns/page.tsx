import { CampaignManagerPage } from "@/features/dashboard/components/workspaces";
import { getDashboardData } from "@/features/dashboard/server/get-dashboard-data";

export default async function CampaignsPage() {
  const data = await getDashboardData();

  return <CampaignManagerPage data={data} />;
}
