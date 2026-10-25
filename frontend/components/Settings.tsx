"use client";

import { useState, useTransition } from "react";
import { updateImageType } from "@/lib/actions";
import { IMAGE_TYPE_LABELS, type ImageType, type User } from "@/lib/types";

const IMAGE_TYPES: { value: ImageType; description: string; icon: string }[] = [
  { value: "photo", description: "実物の写真で覚える", icon: "📷" },
  { value: "illustration", description: "やさしいイラストで覚える", icon: "🎨" },
];

// 設定（GET /me で表示し、PUT /me/settings で画像タイプを保存する）
export function Settings({ user }: { user: User }) {
  const [imageType, setImageType] = useState(user.image_type);
  const [message, setMessage] = useState<{ ok: boolean; text: string } | null>(null);
  const [pending, startTransition] = useTransition();

  function select(value: ImageType) {
    const previous = imageType;
    setImageType(value);
    setMessage(null);
    startTransition(async () => {
      try {
        await updateImageType(value);
        setMessage({ ok: true, text: "保存しました" });
      } catch {
        setImageType(previous); // 保存できなかったら元に戻す
        setMessage({ ok: false, text: "保存できませんでした。もう一度お試しください" });
      }
    });
  }

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <h1 className="text-2xl font-bold">設定</h1>

      <section className="mt-6 rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm">
        <h2 className="font-semibold">画像タイプ</h2>
        <p className="mt-1 text-sm text-zinc-500">単語学習で表示する画像を選べます</p>
        <div className="mt-4 grid grid-cols-2 gap-3">
          {IMAGE_TYPES.map((type) => {
            const selected = imageType === type.value;
            return (
              <button
                key={type.value}
                type="button"
                onClick={() => select(type.value)}
                disabled={pending}
                aria-pressed={selected}
                className={`rounded-2xl border p-4 text-center transition disabled:opacity-60 ${selected ? "border-orange-400 bg-orange-50" : "border-zinc-200 hover:border-orange-200"}`}
              >
                <span className="text-3xl">{type.icon}</span>
                <p className={`mt-2 font-semibold ${selected ? "text-orange-600" : ""}`}>{IMAGE_TYPE_LABELS[type.value]}</p>
                <p className="mt-1 text-xs text-zinc-500">{type.description}</p>
              </button>
            );
          })}
        </div>
        <p className={`mt-3 h-5 text-sm ${message?.ok ? "text-emerald-600" : "text-rose-600"}`} role="status">
          {pending ? "保存中…" : message?.text}
        </p>
      </section>

      <section className="mt-4 rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm">
        <h2 className="font-semibold">アカウント</h2>
        <dl className="mt-3 grid grid-cols-[7rem_1fr] gap-y-2 text-sm">
          <dt className="text-zinc-500">名前</dt>
          <dd>{user.name}</dd>
          <dt className="text-zinc-500">メールアドレス</dt>
          <dd className="break-all">{user.email}</dd>
        </dl>
      </section>
    </main>
  );
}
