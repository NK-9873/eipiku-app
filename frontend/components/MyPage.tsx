"use client";

import Link from "next/link";
import { getCourseWords, getWords, MASTERY_LABELS, type Mastery } from "@/lib/learning";
import { isInterrupted } from "@/lib/review";
import { useLearner } from "./useLearner";

const STAT_STYLES: Record<Mastery, string> = {
  weak: "text-rose-500",
  unsure: "text-amber-500",
  mastered: "text-emerald-500",
};

const card = "rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm";

export function MyPage() {
  const learner = useLearner();
  if (!learner) {
    return <p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>;
  }

  const { user, wordProgresses, courseProgresses } = learner;
  const courses = learner.courses.map((course) => {
    const total = getCourseWords(course.id).length;
    const progress = courseProgresses.find((cp) => cp.course_id === course.id);
    return { ...course, total, position: Math.min(progress?.current_position ?? 0, total), updatedAt: progress?.updated_at ?? "" };
  });
  // 中断しているコースのうち、最後に学習したもの
  const resumeCourse = courses
    .filter((c) => isInterrupted(c.id, courseProgresses))
    .sort((a, b) => b.updatedAt.localeCompare(a.updatedAt))[0];
  const countOf = (mastery: Mastery) => wordProgresses.filter((wp) => wp.mastery === mastery).length;
  const reviewCount = countOf("weak") + countOf("unsure");

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <h1 className="text-2xl font-bold">マイページ</h1>
      <p className="mt-1 text-sm text-zinc-500">{user.name} さん、今日もがんばりましょう</p>

      {resumeCourse && (
        <Link
          href={`/courses/${resumeCourse.id}`}
          className="mt-6 block rounded-2xl bg-orange-500 p-5 text-white shadow-sm transition hover:bg-orange-600 active:scale-[0.99]"
        >
          <p className="text-sm text-orange-100">続きから学習</p>
          <p className="mt-1 text-xl font-bold">{resumeCourse.name}</p>
          <p className="mt-1 text-sm text-orange-100">
            {resumeCourse.position} / {resumeCourse.total} 語まで学習済み ・ {resumeCourse.position + 1} 語目から再開 →
          </p>
        </Link>
      )}

      <section className={`${card} mt-4`}>
        <h2 className="font-semibold">学習状況</h2>
        <p className="mt-3">
          <span className="text-3xl font-bold">{wordProgresses.length}</span>
          <span className="ml-1 text-zinc-500">/ {getWords().length} 語 学習済み</span>
        </p>
        <div className="mt-4 grid grid-cols-3 gap-2 text-center">
          {(["mastered", "unsure", "weak"] as const).map((mastery) => (
            <div key={mastery} className="rounded-xl bg-zinc-50 py-3">
              <p className="text-xs text-zinc-500">{MASTERY_LABELS[mastery]}</p>
              <p className={`text-2xl font-bold ${STAT_STYLES[mastery]}`}>{countOf(mastery)}</p>
            </div>
          ))}
        </div>
      </section>

      <section className={`${card} mt-4`}>
        <div className="flex items-center justify-between gap-3">
          <div>
            <h2 className="font-semibold">復習</h2>
            <p className="mt-1 text-sm text-zinc-500">
              苦手 {countOf("weak")} 語 ・ 微妙 {countOf("unsure")} 語
            </p>
          </div>
          {reviewCount > 0 ? (
            <Link href="/review" className="shrink-0 rounded-xl bg-orange-500 px-4 py-2 font-semibold text-white hover:bg-orange-600">
              復習開始
            </Link>
          ) : (
            <span className="shrink-0 text-sm text-zinc-400">復習する単語なし</span>
          )}
        </div>
      </section>

      <section className={`${card} mt-4`}>
        <h2 className="font-semibold">コース別の進み具合</h2>
        <ul className="mt-3 flex flex-col gap-4">
          {courses.map((course) => (
            <li key={course.id}>
              <div className="flex items-center justify-between gap-3 text-sm">
                <Link href={`/courses/${course.id}`} className="font-medium hover:text-orange-500">
                  {course.name}
                </Link>
                <span className="text-zinc-500">
                  {course.position >= course.total ? "完了" : `${course.position} / ${course.total} 語`}
                </span>
              </div>
              <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-zinc-200">
                <div
                  className={`h-full ${course.position >= course.total ? "bg-emerald-400" : "bg-orange-400"}`}
                  style={{ width: `${course.total ? (course.position / course.total) * 100 : 0}%` }}
                />
              </div>
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}
