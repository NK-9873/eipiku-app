"use client";

import { useEffect } from "react";

// API に繋がらないときなどに表示する
export default function MainError({ error, retry }: { error: unknown; retry: () => void }) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-10">
      <div className="rounded-2xl border border-zinc-200 bg-white p-6 text-center shadow-sm">
        <p className="text-lg font-semibold">データを読み込めませんでした</p>
        <p className="mt-2 text-sm text-zinc-500">
          {error instanceof Error ? error.message : "時間をおいてもう一度お試しください"}
        </p>
        <button
          type="button"
          onClick={() => retry()}
          className="mt-6 rounded-xl bg-orange-500 px-4 py-2 font-semibold text-white transition hover:bg-orange-600 active:scale-95"
        >
          もう一度読み込む
        </button>
      </div>
    </main>
  );
}
