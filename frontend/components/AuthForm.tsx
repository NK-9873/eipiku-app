"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import type { FormEvent } from "react";

const inputClass =
  "w-full rounded-xl border border-zinc-300 bg-white px-3 py-2.5 outline-none focus:border-orange-400 focus:ring-2 focus:ring-orange-100";

// ログイン・新規登録の画面。今は見た目と画面遷移だけで、入力内容はどこにも送らない
export function AuthForm({ mode }: { mode: "login" | "signup" }) {
  const router = useRouter();
  const isLogin = mode === "login";

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    // TODO: user-service の API ができたら、ここで認証・登録する
    router.push("/mypage");
  }

  return (
    <main className="mx-auto flex w-full max-w-sm flex-1 flex-col justify-center px-4 py-10">
      <div className="text-center">
        <p className="text-4xl font-bold text-orange-500">えいピク</p>
        <p className="mt-1 text-sm text-zinc-500">絵で覚える英単語帳</p>
      </div>

      <form onSubmit={handleSubmit} className="mt-8 flex flex-col gap-4 rounded-3xl border border-zinc-200 bg-white p-6 shadow-sm">
        <h1 className="text-xl font-bold">{isLogin ? "ログイン" : "新規登録"}</h1>
        {!isLogin && (
          <label className="block">
            <span className="mb-1 block text-sm font-medium text-zinc-700">名前</span>
            <input name="name" required maxLength={50} autoComplete="name" className={inputClass} />
          </label>
        )}
        <label className="block">
          <span className="mb-1 block text-sm font-medium text-zinc-700">メールアドレス</span>
          <input name="email" type="email" required autoComplete="email" className={inputClass} />
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
        <button
          type="submit"
          className="mt-2 rounded-2xl bg-orange-500 py-3 font-semibold text-white transition hover:bg-orange-600 active:scale-95"
        >
          {isLogin ? "ログイン" : "登録する"}
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
