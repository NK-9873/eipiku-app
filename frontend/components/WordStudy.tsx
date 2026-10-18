"use client";

import Link from "next/link";
import { useState } from "react";
import type { Course, Mastery, Word } from "@/lib/learning";

const MASTERY_BUTTONS: { value: Mastery; label: string; className: string }[] = [
  { value: "weak", label: "苦手", className: "bg-rose-500 hover:bg-rose-600" },
  { value: "unsure", label: "微妙", className: "bg-amber-500 hover:bg-amber-600" },
  { value: "mastered", label: "覚えた", className: "bg-emerald-500 hover:bg-emerald-600" },
];

type Props = {
  course: Course;
  words: Word[];
};

export function WordStudy({ course, words }: Props) {
  const [index, setIndex] = useState(0);
  const [revealed, setRevealed] = useState(false);
  // API ができたら word_progresses に保存する。今は画面内で保持するだけ。
  const [results, setResults] = useState<Record<number, Mastery>>({});

  const word = words[index];
  const finished = index >= words.length;

  function rate(mastery: Mastery) {
    setResults((prev) => ({ ...prev, [word.id]: mastery }));
    setRevealed(false);
    setIndex((i) => i + 1);
  }

  return (
    <main className="mx-auto flex w-full max-w-md flex-1 flex-col px-4 py-6">
      <header className="flex items-center justify-between">
        <h1 className="font-semibold">{course.name}</h1>
        <Link href="/" className="rounded-full px-3 py-1 text-sm text-zinc-500 hover:bg-zinc-100">
          {finished ? "コース一覧へ" : "途中終了"}
        </Link>
      </header>

      {finished ? (
        <Summary words={words} results={results} />
      ) : (
        <>
          <p className="mt-2 text-sm text-zinc-500">
            {index + 1} / {words.length}
          </p>
          <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-zinc-200">
            <div
              className="h-full bg-orange-400 transition-all"
              style={{ width: `${(index / words.length) * 100}%` }}
            />
          </div>

          <button
            type="button"
            onClick={() => setRevealed(true)}
            disabled={revealed}
            className="mt-6 flex flex-col items-center rounded-3xl border border-zinc-200 bg-white p-6 text-center shadow-sm"
          >
            <ImagePlaceholder word={word} />
            <p className="mt-6 text-4xl font-bold tracking-wide">{word.english}</p>

            {revealed ? (
              <div className="mt-4">
                <p className="text-2xl font-semibold text-orange-500">{word.japanese}</p>
                {word.example && <p className="mt-3 text-zinc-600">{word.example}</p>}
              </div>
            ) : (
              <p className="mt-6 text-sm text-zinc-400">タップして意味を見る</p>
            )}
          </button>

          {revealed && (
            <div className="mt-6 grid grid-cols-3 gap-3">
              {MASTERY_BUTTONS.map((b) => (
                <button
                  key={b.value}
                  type="button"
                  onClick={() => rate(b.value)}
                  className={`rounded-2xl py-4 font-semibold text-white transition active:scale-95 ${b.className}`}
                >
                  {b.label}
                </button>
              ))}
            </div>
          )}
        </>
      )}
    </main>
  );
}

function ImagePlaceholder({ word }: { word: Word }) {
  // 画像ができたら word.photo_path / word.illustration_path を表示する
  return (
    <div className="flex aspect-square w-full max-w-64 flex-col items-center justify-center rounded-2xl bg-zinc-100 text-zinc-400">
      <svg viewBox="0 0 24 24" className="h-12 w-12" fill="none" stroke="currentColor" strokeWidth={1.5}>
        <rect x="3" y="4" width="18" height="16" rx="2" />
        <circle cx="9" cy="10" r="2" />
        <path d="m21 16-5-5-9 9" />
      </svg>
      <span className="mt-2 text-xs">{word.english} の画像</span>
    </div>
  );
}

function Summary({ words, results }: { words: Word[]; results: Record<number, Mastery> }) {
  return (
    <div className="mt-10 flex flex-col items-center text-center">
      <p className="text-2xl font-bold">おつかれさま！</p>
      <p className="mt-1 text-zinc-500">{words.length} 語を学習しました</p>

      <div className="mt-8 grid w-full grid-cols-3 gap-3">
        {MASTERY_BUTTONS.map((b) => (
          <div key={b.value} className="rounded-2xl border border-zinc-200 bg-white py-4">
            <p className="text-sm text-zinc-500">{b.label}</p>
            <p className="text-2xl font-bold">
              {words.filter((w) => results[w.id] === b.value).length}
            </p>
          </div>
        ))}
      </div>

      <Link
        href="/"
        className="mt-10 w-full rounded-2xl bg-orange-500 py-4 font-semibold text-white hover:bg-orange-600"
      >
        コース一覧に戻る
      </Link>
    </div>
  );
}
