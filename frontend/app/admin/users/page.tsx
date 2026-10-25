import { Suspense } from "react";
import { UserList } from "@/components/admin/UserViews";
import { Loading } from "@/components/admin/ui";
import { getUsers, requireAdmin } from "@/lib/admin";

export default function AdminUsersPage() {
  return (
    <Suspense fallback={<Loading />}>
      <Users />
    </Suspense>
  );
}

async function Users() {
  const [me, users] = await Promise.all([requireAdmin(), getUsers()]);
  return <UserList users={users} meId={me.id} />;
}
