"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useState } from "react";
import { getCourse, getCourses } from "@/lib/learning";
import { isInterrupted, pickReviewWords, REVIEW_TARGETS, type ReviewTarget } from "@/lib/review";
import { useLearner, type LearnerData } from "./useLearner";

// /review            … マイページから。コースと対象を選び、終わったらマイページへ戻る
// /review?course=1   … 単語学習から。コースは固定で対象だけ選び、終わったら学習の続きへ戻る
export function ReviewSelect() {
  const learner = useLearner();
  const fixedCourseId = Number(useSearchParams().get("course")) || null;

  if (!learner) {
    return <p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>;
  }
  return <ReviewSelectForm learner={learner} fixedCourseId={fixedCourseId} />;
}

function ReviewSelectForm({ learner, fixedCourseId }: { learner: LearnerData; fixedCourseId: number | null }) {
  const fixedCourse = fixedCourseId ? getCourse(fixedCourseId) : undefined;
  // 中断中のコースを先に出す（sort は安定なので、それ以外は元の順番のまま）
  const courses = [...getCourses()].sort(
    (a, b) =>
      Number(isInterrupted(b.id, learner.courseProgresses)) - Number(isInterrupted(a.id, learner.courseProgresses)),
  );

  const [courseId, setCourseId] = useState<number | "all">(fixedCourse?.id ?? "all");
  const [target, setTarget] = useState<ReviewTarget>("both");

  const count = (c: number | "all", t: ReviewTarget) => pickReviewWords(c, t, learner.wordProgresses).length;
  const selectedCount = count(courseId, target);
  const back = fixedCourse ? "course" : "mypage";
  const backHref = fixedCourse ? `/courses/${fixedCourse.id}` : "/mypage";

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <Link href={backHref} className="mb-3 inline-block text-sm text-zinc-500 hover:text-orange-500">
        ← {fixedCourse ? "学習に戻る" : "マイページ"}
      </Link>
      <h1 className="text-2xl font-bold">復習する単語を選ぶ</h1>

      <section className="mt-6">
        <h2 className="mb-2 font-semibold">コース</h2>
        {fixedCourse ? (
          <p className="rounded-2xl border border-orange-300 bg-white p-4 font-semibold">{fixedCourse.name}</p>
        ) : (
          <div className="flex flex-col gap-2">
            <Choice selected={courseId === "all"} onClick={() => setCourseId("all")} label="すべてのコース" count={count("all", target)} />
            {courses.map((course) => (
              <Choice
                key={course.id}
                selected={courseId === course.id}
                onClick={() => setCourseId(course.id)}
                label={course.name}
                badge={isInterrupted(course.id, learner.courseProgresses) ? "中断中" : undefined}
                count={count(course.id, target)}
              />
            ))}
          </div>
        )}
      </section>

      <section className="mt-6">
        <h2 className="mb-2 font-semibold">対象</h2>
        <div className="grid grid-cols-3 gap-2">
          {REVIEW_TARGETS.map((t) => (
            <button
              key={t.value}
              type="button"
              onClick={() => setTarget(t.value)}
              className={`rounded-2xl border py-3 transition ${target === t.value ? "border-orange-400 bg-orange-50 font-semibold text-orange-600" : "border-zinc-200 bg-white text-zinc-700 hover:border-orange-200"}`}
            >
              {t.label}
              <span className="block text-xs font-normal text-zinc-500">{count(courseId, t.value)} 語</span>
            </button>
          ))}
        </div>
      </section>

      {selectedCount > 0 ? (
        <Link
          href={`/review/session?course=${courseId}&target=${target}&back=${back}`}
          className="mt-8 block w-full rounded-2xl bg-orange-500 py-4 text-center font-semibold text-white hover:bg-orange-600"
        >
          {selectedCount} 語を復習する
        </Link>
      ) : (
        <p className="mt-8 rounded-2xl bg-zinc-100 py-4 text-center text-sm text-zinc-500">
          復習する単語はありません
        </p>
      )}
    </main>
  );
}

function Choice({
  selected,
  onClick,
  label,
  badge,
  count,
}: {
  selected: boolean;
  onClick: () => void;
  label: string;
  badge?: string;
  count: number;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`flex items-center justify-between rounded-2xl border p-4 text-left transition ${selected ? "border-orange-400 bg-orange-50" : "border-zinc-200 bg-white hover:border-orange-200"}`}
    >
      <span className="flex items-center gap-2">
        <span className={selected ? "font-semibold text-orange-600" : ""}>{label}</span>
        {badge && <span className="rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-600">{badge}</span>}
      </span>
      <span className="text-sm text-zinc-500">{count} 語</span>
    </button>
  );
}
