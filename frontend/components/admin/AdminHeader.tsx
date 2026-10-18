"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const NAV_ITEMS = [
  { href: "/admin/courses", label: "コース" },
  { href: "/admin/words", label: "単語" },
  { href: "/admin/users", label: "ユーザー" },
];

export function AdminHeader() {
  const pathname = usePathname();

  return (
    <header className="border-b border-orange-100 bg-white/80">
      <div className="mx-auto flex w-full max-w-3xl flex-wrap items-center justify-between gap-2 px-4 py-3">
        <Link href="/admin" className="flex items-baseline gap-2">
          <span className="text-xl font-bold text-orange-500">えいピク</span>
          <span className="text-sm font-semibold text-zinc-500">管理画面</span>
        </Link>
        <nav className="flex items-center gap-1 text-sm">
          {NAV_ITEMS.map((item) => {
            const active = pathname.startsWith(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`rounded-full px-3 py-1 ${active ? "bg-orange-100 font-semibold text-orange-600" : "text-zinc-600 hover:bg-zinc-100"}`}
              >
                {item.label}
              </Link>
            );
          })}
          <Link href="/" className="ml-2 rounded-full px-3 py-1 text-zinc-400 hover:bg-zinc-100">
            アプリへ
          </Link>
        </nav>
      </div>
    </header>
  );
}
