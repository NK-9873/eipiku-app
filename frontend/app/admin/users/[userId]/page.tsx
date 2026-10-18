import { Suspense } from "react";
import { UserDetail } from "@/components/admin/UserViews";
import { Loading } from "@/components/admin/ui";

export default function AdminUserDetailPage() {
  return (
    <Suspense fallback={<Loading />}>
      <UserDetail />
    </Suspense>
  );
}
