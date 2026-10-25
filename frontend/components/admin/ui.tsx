// 管理画面で共通して使う見た目の部品
import Link from "next/link";
import type { ReactNode } from "react";

export const cardClass = "rounded-2xl border border-zinc-200 bg-white shadow-sm";
export const inputClass =
  "w-full rounded-xl border border-zinc-300 bg-white px-3 py-2 outline-none focus:border-orange-400 focus:ring-2 focus:ring-orange-100";
export const primaryButtonClass =
  "rounded-xl bg-orange-500 px-4 py-2 font-semibold text-white transition hover:bg-orange-600 active:scale-95";
export const secondaryButtonClass =
  "rounded-xl border border-zinc-300 bg-white px-4 py-2 font-semibold text-zinc-700 transition hover:bg-zinc-50 active:scale-95";
export const smallButtonClass =
  "rounded-lg border border-zinc-200 px-2.5 py-1 text-sm text-zinc-600 transition hover:bg-zinc-50 disabled:opacity-30";
export const dangerButtonClass =
  "rounded-lg px-2.5 py-1 text-sm text-rose-600 transition hover:bg-rose-50";

export function PageHeader({ title, action }: { title: string; action?: ReactNode }) {
  return (
    <div className="mb-5 flex items-center justify-between gap-3">
      <h1 className="text-2xl font-bold">{title}</h1>
      {action}
    </div>
  );
}

export function BackLink({ href, children }: { href: string; children: ReactNode }) {
  return (
    <Link href={href} className="mb-3 inline-block text-sm text-zinc-500 hover:text-orange-500">
      ← {children}
    </Link>
  );
}

export function Field({ label, required, children }: { label: string; required?: boolean; children: ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1 block text-sm font-medium text-zinc-700">
        {label}
        {required && <span className="ml-1 text-rose-500">*</span>}
      </span>
      {children}
    </label>
  );
}

export function NotFoundMessage({ backHref, children }: { backHref: string; children: ReactNode }) {
  return (
    <div className={`${cardClass} p-8 text-center`}>
      <p className="text-zinc-500">{children}</p>
      <Link href={backHref} className={`${secondaryButtonClass} mt-4 inline-block`}>
        一覧に戻る
      </Link>
    </div>
  );
}

export function Loading() {
  return <p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>;
}

// API のエラーを画面に出す（details は項目ごとの理由）
export function ErrorBox({ error, details }: { error?: string; details?: string[] }) {
  if (!error) return null;
  return (
    <div role="alert" className="rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
      <p className="font-semibold">{error}</p>
      {details && details.length > 0 && (
        <ul className="mt-1 list-disc pl-5">
          {details.map((d) => (
            <li key={d}>{d}</li>
          ))}
        </ul>
      )}
    </div>
  );
}

// 保存・削除が終わったことを知らせる
export function SavedNotice({ children }: { children: ReactNode }) {
  return (
    <p role="status" className="mb-4 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {children}
    </p>
  );
}
