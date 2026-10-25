"use client";

import Image from "next/image";
import { useState } from "react";
import { IMAGE_TYPE_LABELS, MASTERY_LABELS, type ImageType, type Mastery, type Word } from "@/lib/types";

export const MASTERY_STYLES: Record<Mastery, string> = {
  weak: "bg-rose-500 hover:bg-rose-600",
  unsure: "bg-amber-500 hover:bg-amber-600",
  mastered: "bg-emerald-500 hover:bg-emerald-600",
};

const MASTERIES: Mastery[] = ["weak", "unsure", "mastered"];

type Props = {
  word: Word;
  imageType: ImageType;
  onRate: (mastery: Mastery) => void;
};

// 画像と英単語 → タップで訳・例文・例文の訳 → 苦手・微妙・覚えた で評価、の1語分。
// 単語が変わったら表示を戻したいので、呼び出し側で key={word.id} を付ける
export function WordCard({ word, imageType, onRate }: Props) {
  const [revealed, setRevealed] = useState(false);

  return (
    <>
      <button
        type="button"
        onClick={() => setRevealed(true)}
        disabled={revealed}
        className="mt-6 flex flex-col items-center rounded-3xl border border-zinc-200 bg-white p-6 text-center shadow-sm"
      >
        <WordImage word={word} imageType={imageType} />
        <p className="mt-6 text-4xl font-bold tracking-wide">{word.english}</p>

        {revealed ? (
          <div className="mt-4">
            <p className="text-2xl font-semibold text-orange-500">{word.japanese}</p>
            {word.example && <p className="mt-3 text-zinc-700">{word.example}</p>}
            {word.example_japanese && <p className="mt-1 text-sm text-zinc-500">{word.example_japanese}</p>}
          </div>
        ) : (
          <p className="mt-6 text-sm text-zinc-400">タップして意味を見る</p>
        )}
      </button>

      {revealed && (
        <div className="mt-6 grid grid-cols-3 gap-3">
          {MASTERIES.map((mastery) => (
            <button
              key={mastery}
              type="button"
              onClick={() => onRate(mastery)}
              className={`rounded-2xl py-4 font-semibold text-white transition active:scale-95 ${MASTERY_STYLES[mastery]}`}
            >
              {MASTERY_LABELS[mastery]}
            </button>
          ))}
        </div>
      )}
    </>
  );
}

// 設定の画像タイプに合わせて写真かイラストを出す。
// URL は learning-service の POST /images が返したもの。登録されていない・読み込めないときはプレースホルダー
function WordImage({ word, imageType }: { word: Word; imageType: ImageType }) {
  const src = imageType === "photo" ? word.photo_path : word.illustration_path;
  const [failedSrc, setFailedSrc] = useState<string | null>(null);

  return (
    <div className="relative flex aspect-square w-full max-w-64 flex-col items-center justify-center overflow-hidden rounded-2xl bg-zinc-100 text-zinc-400">
      {src && failedSrc !== src ? (
        <Image
          src={src}
          alt={word.english}
          fill
          sizes="256px"
          unoptimized
          className="object-cover"
          onError={() => setFailedSrc(src)}
        />
      ) : (
        <>
          <svg viewBox="0 0 24 24" className="h-12 w-12" fill="none" stroke="currentColor" strokeWidth={1.5}>
            <rect x="3" y="4" width="18" height="16" rx="2" />
            <circle cx="9" cy="10" r="2" />
            <path d="m21 16-5-5-9 9" />
          </svg>
          <span className="mt-2 text-xs">
            {word.english} の{IMAGE_TYPE_LABELS[imageType]}
          </span>
        </>
      )}
    </div>
  );
}
