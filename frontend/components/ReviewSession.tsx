"use client";

import Link from "next/link";
import { useState } from "react";
import { saveWordProgress } from "@/lib/actions";
import { MASTERY_LABELS, type ImageType, type Mastery, type ReviewWord } from "@/lib/types";
import { WordCard } from "./WordCard";

const MASTERIES: Mastery[] = ["weak", "unsure", "mastered"];

type Props = {
  words: ReviewWord[];
  targetLabel: string;
  imageType: ImageType;
  backHref: string;
  backLabel: string;
};

export function ReviewSession({ words, targetLabel, imageType, backHref, backLabel }: Props) {
  const [index, setIndex] = useState(0);
  const [results, setResults] = useState<Mastery[]>([]);
  const [saveError, setSaveError] = useState(false);

  const word = words[index];
  const finished = index >= words.length;

  function rate(mastery: Mastery) {
    // 待たずに次の単語へ進み、習得度の更新は裏で行う
    saveWordProgress(word.id, mastery).catch(() => setSaveError(true));
    setResults((prev) => [...prev, mastery]);
    setIndex((i) => i + 1);
  }

  return (
    <main className="mx-auto flex w-full max-w-md flex-1 flex-col px-4 py-6">
      <header className="flex items-center justify-between">
        <h1 className="font-semibold">
          復習 <span className="text-sm font-normal text-zinc-500">（{targetLabel}）</span>
        </h1>
        {!finished && (
          <Link href={backHref} className="rounded-full px-3 py-1 text-sm text-zinc-500 hover:bg-zinc-100">
            途中終了
          </Link>
        )}
      </header>

      {saveError && (
        <p role="alert" className="mt-3 rounded-xl bg-rose-50 px-4 py-2 text-sm text-rose-700">
          習得度を保存できませんでした。通信状況を確かめて、もう一度お試しください。
        </p>
      )}

      {words.length === 0 ? (
        <div className="mt-10 text-center">
          <p className="text-zinc-500">復習する単語はありません</p>
          <Link href={backHref} className="mt-6 block w-full rounded-2xl bg-orange-500 py-4 font-semibold text-white hover:bg-orange-600">
            {backLabel}
          </Link>
        </div>
      ) : finished ? (
        <div className="mt-10 flex flex-col items-center text-center">
          <p className="text-2xl font-bold">復習おわり！</p>
          <p className="mt-1 text-zinc-500">
            {words.length} 語のうち {results.filter((r) => r === "mastered").length} 語を覚えました
          </p>
          <div className="mt-8 grid w-full grid-cols-3 gap-3">
            {MASTERIES.map((mastery) => (
              <div key={mastery} className="rounded-2xl border border-zinc-200 bg-white py-4">
                <p className="text-sm text-zinc-500">{MASTERY_LABELS[mastery]}</p>
                <p className="text-2xl font-bold">{results.filter((r) => r === mastery).length}</p>
              </div>
            ))}
          </div>
          <Link href={backHref} className="mt-10 w-full rounded-2xl bg-orange-500 py-4 font-semibold text-white hover:bg-orange-600">
            {backLabel}
          </Link>
        </div>
      ) : (
        <>
          <p className="mt-2 text-sm text-zinc-500">
            {index + 1} / {words.length}
            <span className="ml-2">前回：{MASTERY_LABELS[word.mastery]}</span>
          </p>
          <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-zinc-200">
            <div className="h-full bg-orange-400 transition-all" style={{ width: `${(index / words.length) * 100}%` }} />
          </div>
          <WordCard key={word.id} word={word} imageType={imageType} onRate={rate} />
        </>
      )}
    </main>
  );
}
