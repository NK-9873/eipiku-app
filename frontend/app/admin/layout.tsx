import type { Metadata } from "next";
import { Suspense } from "react";
import { AdminHeader } from "@/components/admin/AdminHeader";
import { Loading } from "@/components/admin/ui";
import { requireAdmin } from "@/lib/admin";

export const metadata: Metadata = {
  title: "えいピク 管理画面",
};

// /admin 以下は管理者だけ。一般ユーザーはトップへ、未ログインはログイン画面へ移す。
// （各 API も管理者でなければ 403 を返すので、ここを通り抜けてもデータは見えない）
export default function AdminLayout({ children }: LayoutProps<"/admin">) {
  return (
    <Suspense fallback={<Loading />}>
      <AdminGate>{children}</AdminGate>
    </Suspense>
  );
}

async function AdminGate({ children }: { children: React.ReactNode }) {
  const me = await requireAdmin();
  return (
    <>
      <AdminHeader userName={me.name} />
      <main className="mx-auto w-full max-w-3xl flex-1 px-4 py-8">{children}</main>
    </>
  );
}
