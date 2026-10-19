import { TokenPlanDashboard } from "@/components/token-plan-dashboard";

// Token plan page. The backing platform is an implementation detail and is
// never named in the UI or the URL.
export default function TokenPlanPage() {
  return <TokenPlanDashboard site="tpai" />;
}
