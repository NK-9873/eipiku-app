import { Suspense } from "react";
import { MyPage } from "@/components/MyPage";
import { getSummary } from "@/lib/progress";
import { getMe } from "@/lib/user";

export default function MyPagePage() {
  return (
    <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
      <MyPageData />
    </Suspense>
  );
}

async function MyPageData() {
  const [me, summary] = await Promise.all([getMe(), getSummary()]);
  return <MyPage userName={me.name} summary={summary} />;
}
