"use client";

import { useState } from "react";
import { IMAGE_TYPE_LABELS, updateImageType, type ImageType } from "@/lib/user";
import { useLearner, type LearnerData } from "./useLearner";

const IMAGE_TYPES: { value: ImageType; description: string; icon: string }[] = [
  { value: "photo", description: "実物の写真で覚える", icon: "📷" },
  { value: "illustration", description: "やさしいイラストで覚える", icon: "🎨" },
];

export function Settings() {
  const learner = useLearner();
  if (!learner) {
    return <p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>;
  }
  return <SettingsForm learner={learner} />;
}

function SettingsForm({ learner }: { learner: LearnerData }) {
  const [imageType, setImageType] = useState(learner.user.image_type);
  const [saved, setSaved] = useState(false);

  async function select(value: ImageType) {
    setImageType(value);
    await updateImageType(value);
    setSaved(true);
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
                aria-pressed={selected}
                className={`rounded-2xl border p-4 text-center transition ${selected ? "border-orange-400 bg-orange-50" : "border-zinc-200 hover:border-orange-200"}`}
              >
                <span className="text-3xl">{type.icon}</span>
                <p className={`mt-2 font-semibold ${selected ? "text-orange-600" : ""}`}>{IMAGE_TYPE_LABELS[type.value]}</p>
                <p className="mt-1 text-xs text-zinc-500">{type.description}</p>
              </button>
            );
          })}
        </div>
        <p className="mt-3 h-5 text-sm text-emerald-600">{saved && "保存しました"}</p>
      </section>

      <section className="mt-4 rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm">
        <h2 className="font-semibold">アカウント</h2>
        <dl className="mt-3 grid grid-cols-[7rem_1fr] gap-y-2 text-sm">
          <dt className="text-zinc-500">名前</dt>
          <dd>{learner.user.name}</dd>
          <dt className="text-zinc-500">メールアドレス</dt>
          <dd className="break-all">{learner.user.email}</dd>
        </dl>
      </section>
    </main>
  );
}
