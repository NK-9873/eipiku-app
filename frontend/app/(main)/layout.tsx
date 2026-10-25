import { Suspense } from "react";
import { AppHeader } from "@/components/AppHeader";
import { getMe } from "@/lib/user";

export default function MainLayout({ children }: LayoutProps<"/">) {
  return (
    <>
      {/* ユーザー名は API から取るので、届くまでは名前なしのヘッダーを出す */}
      <Suspense fallback={<HeaderPlaceholder />}>
        <HeaderWithUser />
      </Suspense>
      {children}
    </>
  );
}

async function HeaderWithUser() {
  const me = await getMe();
  return <AppHeader userName={me.name} />;
}

function HeaderPlaceholder() {
  return (
    <header className="border-b border-orange-100 bg-white/80">
      <div className="mx-auto w-full max-w-md px-4 py-3">
        <span className="text-xl font-bold text-orange-500">えいピク</span>
        <div className="mt-2 h-7" />
      </div>
    </header>
  );
}
