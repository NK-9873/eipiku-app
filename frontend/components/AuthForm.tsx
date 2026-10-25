"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useActionState } from "react";
import { login, signup } from "@/lib/actions";
import type { FormState } from "@/lib/types";

const inputClass =
  "w-full rounded-xl border border-zinc-300 bg-white px-3 py-2.5 outline-none focus:border-orange-400 focus:ring-2 focus:ring-orange-100";

// ログイン（user-service の POST /login）・新規登録（POST /signup）の画面
export function AuthForm({ mode }: { mode: "login" | "signup" }) {
  const isLogin = mode === "login";
  const [state, formAction, pending] = useActionState<FormState, FormData>(isLogin ? login : signup, {});
  const registered = useSearchParams().get("registered") === "1";

  return (
    <main className="mx-auto flex w-full max-w-sm flex-1 flex-col justify-center px-4 py-10">
      <div className="text-center">
        <p className="text-4xl font-bold text-orange-500">えいピク</p>
        <p className="mt-1 text-sm text-zinc-500">絵で覚える英単語帳</p>
      </div>

      {isLogin && registered && !state.error && (
        <p className="mt-6 rounded-xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
          登録しました。メールアドレスとパスワードでログインしてください。
        </p>
      )}

      <form action={formAction} className="mt-6 flex flex-col gap-4 rounded-3xl border border-zinc-200 bg-white p-6 shadow-sm">
        <h1 className="text-xl font-bold">{isLogin ? "ログイン" : "新規登録"}</h1>
        {!isLogin && (
          <label className="block">
            <span className="mb-1 block text-sm font-medium text-zinc-700">名前</span>
            <input name="name" required maxLength={50} autoComplete="name" defaultValue={state.values?.name} className={inputClass} />
          </label>
        )}
        <label className="block">
          <span className="mb-1 block text-sm font-medium text-zinc-700">メールアドレス</span>
          <input
            name="email"
            type="email"
            required
            autoComplete="email"
            defaultValue={state.values?.email}
            className={inputClass}
          />
        </label>
        <label className="block">
          <span className="mb-1 block text-sm font-medium text-zinc-700">パスワード</span>
          <input
            name="password"
            type="password"
            required
            minLength={isLogin ? undefined : 8}
            autoComplete={isLogin ? "current-password" : "new-password"}
            className={inputClass}
          />
          {!isLogin && <span className="mt-1 block text-xs text-zinc-400">8文字以上</span>}
        </label>

        {state.error && (
          <div role="alert" className="rounded-xl bg-rose-50 px-4 py-3 text-sm text-rose-700">
            <p>{state.error}</p>
            {state.details && state.details.length > 0 && (
              <ul className="mt-1 list-disc pl-5">
                {state.details.map((d) => (
                  <li key={d}>{d}</li>
                ))}
              </ul>
            )}
          </div>
        )}

        <button
          type="submit"
          disabled={pending}
          className="mt-2 rounded-2xl bg-orange-500 py-3 font-semibold text-white transition hover:bg-orange-600 active:scale-95 disabled:opacity-60"
        >
          {pending ? "送信中…" : isLogin ? "ログイン" : "登録する"}
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-zinc-500">
        {isLogin ? "はじめての方は " : "アカウントをお持ちの方は "}
        <Link href={isLogin ? "/signup" : "/login"} className="font-semibold text-orange-500 hover:underline">
          {isLogin ? "新規登録" : "ログイン"}
        </Link>
      </p>
    </main>
  );
}
