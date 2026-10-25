import { Suspense } from "react";
import { Settings } from "@/components/Settings";
import { getMe } from "@/lib/user";

export default function SettingsPage() {
  return (
    <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
      <SettingsData />
    </Suspense>
  );
}

async function SettingsData() {
  return <Settings user={await getMe()} />;
}
