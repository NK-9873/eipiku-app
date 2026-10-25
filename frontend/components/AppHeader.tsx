"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { logout } from "@/lib/actions";

const NAV_ITEMS = [
  { href: "/mypage", label: "マイページ" },
  { href: "/", label: "コース" },
  { href: "/settings", label: "設定" },
];

function isActive(pathname: string, href: string) {
  // コース一覧は単語学習（/courses/...）の間も選択中にする
  return href === "/" ? pathname === "/" || pathname.startsWith("/courses") : pathname.startsWith(href);
}

// userName が null のあいだ（ユーザー情報の読み込み中）は名前を出さない
export function AppHeader({ userName }: { userName: string | null }) {
  const pathname = usePathname();

  return (
    <header className="border-b border-orange-100 bg-white/80">
      <div className="mx-auto w-full max-w-md px-4 py-3">
        <div className="flex items-center justify-between gap-2">
          <Link href="/mypage" className="text-xl font-bold text-orange-500">
            えいピク
          </Link>
          <span className="truncate text-sm text-zinc-500">{userName && `${userName} さん`}</span>
        </div>
        <nav className="mt-2 flex items-center gap-0.5 text-sm">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`rounded-full px-2.5 py-1 ${isActive(pathname, item.href) ? "bg-orange-100 font-semibold text-orange-600" : "text-zinc-600 hover:bg-zinc-100"}`}
            >
              {item.label}
            </Link>
          ))}
          <form action={logout} className="ml-auto">
            <button type="submit" className="rounded-full px-2.5 py-1 text-zinc-400 hover:bg-zinc-100">
              ログアウト
            </button>
          </form>
        </nav>
      </div>
    </header>
  );
}
