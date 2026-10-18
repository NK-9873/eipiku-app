import type { Metadata } from "next";
import { Suspense } from "react";
import { AdminHeader } from "@/components/admin/AdminHeader";
import { AdminStoreProvider } from "@/components/admin/AdminStore";

export const metadata: Metadata = {
  title: "えいピク 管理画面",
};

// TODO: ログイン機能ができたら、管理者だけが入れるようにする
export default function AdminLayout({ children }: LayoutProps<"/admin">) {
  return (
    <AdminStoreProvider>
      <Suspense>
        <AdminHeader />
      </Suspense>
      <main className="mx-auto w-full max-w-3xl flex-1 px-4 py-8">{children}</main>
    </AdminStoreProvider>
  );
}
