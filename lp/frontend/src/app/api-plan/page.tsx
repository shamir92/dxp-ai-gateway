import { notFound } from "next/navigation";

// Hidden until API plans launch. The page is kept in git so it can be restored
// by reverting this file; the nav entry in public-shell/app-shell is already inert.
export default function ApiPlanPage() {
  notFound();
}
