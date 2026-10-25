"use client";

import Link from "next/link";
import { useState } from "react";
import { recordStudy } from "@/lib/actions";
import { MASTERY_LABELS, type Course, type ImageType, type Mastery, type Word } from "@/lib/types";
import { WordCard } from "./WordCard";

const MASTERIES: Mastery[] = ["weak", "unsure", "mastered"];

type Props = {
  course: Course;
  words: Word[];
  // 中断していればその位置から（サーバーで GET /me/course-progresses から決める）
  startIndex: number;
  imageType: ImageType;
};

export function WordStudy({ course, words, startIndex, imageType }: Props) {
  // 何語目を出しているか（0 始まり）
  const [index, setIndex] = useState(startIndex);
  // 今回の学習で付けた評価（終了画面の集計用）
  const [results, setResults] = useState<Record<number, Mastery>>({});
  const [saveError, setSaveError] = useState(false);

  const word = words[index];
  const finished = index >= words.length;

  function rate(mastery: Mastery) {
    const next = index + 1;
    // 待たずに次の単語へ進み、保存は裏で行う。
    // 1語ごとに中断位置も保存するので、途中終了や画面を閉じたときもそこから再開できる
    recordStudy({ courseId: course.id, wordId: word.id, mastery, nextPosition: next }).catch(() => setSaveError(true));
    setResults((prev) => ({ ...prev, [word.id]: mastery }));
    setIndex(next);
  }

  return (
    <main className="mx-auto flex w-full max-w-md flex-1 flex-col px-4 py-6">
      <header className="flex items-center justify-between gap-2">
        <h1 className="font-semibold">{course.name}</h1>
        <div className="flex items-center gap-1">
          {!finished && (
            <Link
              href={`/review?course=${course.id}`}
              className="rounded-full border border-orange-200 px-3 py-1 text-sm text-orange-600 hover:bg-orange-50"
            >
              復習へ
            </Link>
          )}
          <Link href="/" className="rounded-full px-3 py-1 text-sm text-zinc-500 hover:bg-zinc-100">
            {finished ? "コース一覧へ" : "途中終了"}
          </Link>
        </div>
      </header>

      {saveError && (
        <p role="alert" className="mt-3 rounded-xl bg-rose-50 px-4 py-2 text-sm text-rose-700">
          学習の記録を保存できませんでした。通信状況を確かめて、もう一度お試しください。
        </p>
      )}

      {words.length === 0 ? (
        <p className="mt-10 text-center text-zinc-500">このコースにはまだ単語がありません</p>
      ) : finished ? (
        <Summary learnedCount={words.length - startIndex} results={results} resumed={startIndex > 0} />
      ) : (
        <>
          <p className="mt-2 text-sm text-zinc-500">
            {index + 1} / {words.length}
            {startIndex > 0 && index === startIndex && <span className="ml-2 text-orange-500">続きから再開</span>}
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

function Summary({
  learnedCount,
  results,
  resumed,
}: {
  learnedCount: number;
  results: Record<number, Mastery>;
  resumed: boolean;
}) {
  const values = Object.values(results);

  return (
    <div className="mt-10 flex flex-col items-center text-center">
      <p className="text-2xl font-bold">おつかれさま！</p>
      <p className="mt-1 text-zinc-500">
        {resumed ? "続きから" : ""}
        {learnedCount} 語を学習しました
      </p>

      <div className="mt-8 grid w-full grid-cols-3 gap-3">
        {MASTERIES.map((mastery) => (
          <div key={mastery} className="rounded-2xl border border-zinc-200 bg-white py-4">
            <p className="text-sm text-zinc-500">{MASTERY_LABELS[mastery]}</p>
            <p className="text-2xl font-bold">{values.filter((v) => v === mastery).length}</p>
          </div>
        ))}
      </div>

      <div className="mt-10 flex w-full flex-col gap-3">
        <Link href="/mypage" className="w-full rounded-2xl bg-orange-500 py-4 font-semibold text-white hover:bg-orange-600">
          マイページへ
        </Link>
        <Link href="/" className="w-full rounded-2xl border border-zinc-300 bg-white py-4 font-semibold text-zinc-700 hover:bg-zinc-50">
          コース一覧に戻る
        </Link>
      </div>
    </div>
  );
}
