"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const NAV_ITEMS = [
  { href: "/mypage", label: "マイページ" },
  { href: "/", label: "コース" },
  { href: "/settings", label: "設定" },
];

function isActive(pathname: string, href: string) {
  // コース一覧は単語学習（/courses/...）の間も選択中にする
  return href === "/" ? pathname === "/" || pathname.startsWith("/courses") : pathname.startsWith(href);
}

export function AppHeader() {
  const pathname = usePathname();

  return (
    <header className="border-b border-orange-100 bg-white/80">
      <div className="mx-auto flex w-full max-w-md items-center justify-between gap-2 px-4 py-3">
        <Link href="/mypage" className="text-xl font-bold text-orange-500">
          えいピク
        </Link>
        <nav className="flex items-center gap-0.5 text-sm">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`rounded-full px-2.5 py-1 ${isActive(pathname, item.href) ? "bg-orange-100 font-semibold text-orange-600" : "text-zinc-600 hover:bg-zinc-100"}`}
            >
              {item.label}
            </Link>
          ))}
          {/* TODO: ログイン機能ができたら、ここでセッションを破棄する */}
          <Link href="/login" className="rounded-full px-2.5 py-1 text-zinc-400 hover:bg-zinc-100">
            ログアウト
          </Link>
        </nav>
      </div>
    </header>
  );
}
