"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { MASTERY_LABELS, type Mastery } from "@/lib/learning";
import { saveWordProgress } from "@/lib/progress";
import { pickReviewWords, REVIEW_TARGETS, type ReviewTarget } from "@/lib/review";
import { CURRENT_USER_ID } from "@/lib/user";
import { useLearner, type LearnerData } from "./useLearner";
import { WordCard } from "./WordCard";

const MASTERIES: Mastery[] = ["weak", "unsure", "mastered"];

// URL: /review/session?course=all|<id>&target=weak|unsure|both&back=mypage|course
export function ReviewSession() {
  const learner = useLearner();
  const params = useSearchParams();

  const courseParam = params.get("course");
  const courseId = courseParam && courseParam !== "all" ? Number(courseParam) : "all";
  const target = (REVIEW_TARGETS.find((t) => t.value === params.get("target"))?.value ?? "both") as ReviewTarget;
  // 単語学習から来たときは学習の続き（中断位置）へ、マイページから来たときはマイページへ戻る
  const backHref = params.get("back") === "course" && courseId !== "all" ? `/courses/${courseId}` : "/mypage";
  const backLabel = backHref === "/mypage" ? "マイページに戻る" : "学習を再開する";

  if (!learner) {
    return <p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>;
  }
  return (
    <Session learner={learner} courseId={courseId} target={target} backHref={backHref} backLabel={backLabel} />
  );
}

type SessionProps = {
  learner: LearnerData;
  courseId: number | "all";
  target: ReviewTarget;
  backHref: string;
  backLabel: string;
};

function Session({ learner, courseId, target, backHref, backLabel }: SessionProps) {
  // 出題する単語は開始時点で決めておく（評価して対象から外れても、今回の復習では最後まで出す）
  const [words] = useState(() =>
    pickReviewWords(
      courseId === "all" ? learner.courses.map((course) => course.id) : [courseId],
      target,
      learner.wordProgresses,
    ),
  );
  const [index, setIndex] = useState(0);
  const [results, setResults] = useState<Mastery[]>([]);

  const word = words[index];
  const finished = index >= words.length;

  function rate(mastery: Mastery) {
    saveWordProgress(CURRENT_USER_ID, word.id, mastery);
    setResults((prev) => [...prev, mastery]);
    setIndex((i) => i + 1);
  }

  return (
    <main className="mx-auto flex w-full max-w-md flex-1 flex-col px-4 py-6">
      <header className="flex items-center justify-between">
        <h1 className="font-semibold">
          復習 <span className="text-sm font-normal text-zinc-500">（{REVIEW_TARGETS.find((t) => t.value === target)?.label}）</span>
        </h1>
        {!finished && (
          <Link href={backHref} className="rounded-full px-3 py-1 text-sm text-zinc-500 hover:bg-zinc-100">
            途中終了
          </Link>
        )}
      </header>

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
            <span className="ml-2">
              前回：{MASTERY_LABELS[learner.wordProgresses.find((wp) => wp.word_id === word.id)!.mastery]}
            </span>
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
