"use client";

import Link from "next/link";
import { useState } from "react";
import { MASTERY_LABELS, type Course, type Mastery, type Word } from "@/lib/learning";
import { saveCourseProgress, saveWordProgress } from "@/lib/progress";
import { isInterrupted } from "@/lib/review";
import { CURRENT_USER_ID } from "@/lib/user";
import { useLearner, type LearnerData } from "./useLearner";
import { WordCard } from "./WordCard";

const MASTERIES: Mastery[] = ["weak", "unsure", "mastered"];

type Props = {
  course: Course;
  words: Word[];
};

export function WordStudy({ course, words }: Props) {
  const learner = useLearner();
  if (!learner) {
    return <p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>;
  }
  return <StudySession course={course} words={words} learner={learner} />;
}

function StudySession({ course, words, learner }: Props & { learner: LearnerData }) {
  // 中断していれば続きから、そうでなければ最初から
  const [startIndex] = useState(() =>
    isInterrupted(course.id, learner.courseProgresses)
      ? learner.courseProgresses.find((cp) => cp.course_id === course.id)!.current_position
      : 0,
  );
  // 何語目を出しているか（0 始まり）
  const [index, setIndex] = useState(startIndex);
  // 今回の学習で付けた評価（終了画面の集計用）
  const [results, setResults] = useState<Record<number, Mastery>>({});

  const word = words[index];
  const finished = index >= words.length;

  function rate(mastery: Mastery) {
    const next = index + 1;
    saveWordProgress(CURRENT_USER_ID, word.id, mastery);
    // 1語ごとに中断位置を保存しておけば、途中終了や画面を閉じたときもそこから再開できる
    saveCourseProgress(CURRENT_USER_ID, course.id, next);
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

      {finished ? (
        <Summary
          learnedCount={words.length - startIndex}
          results={results}
          resumed={startIndex > 0}
        />
      ) : (
        <>
          <p className="mt-2 text-sm text-zinc-500">
            {index + 1} / {words.length}
            {startIndex > 0 && index === startIndex && <span className="ml-2 text-orange-500">続きから再開</span>}
          </p>
          <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-zinc-200">
            <div className="h-full bg-orange-400 transition-all" style={{ width: `${(index / words.length) * 100}%` }} />
          </div>
          <WordCard key={word.id} word={word} imageType={learner.user.image_type} onRate={rate} />
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
