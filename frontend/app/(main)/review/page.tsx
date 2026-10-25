import { notFound } from "next/navigation";
import { Suspense } from "react";
import { ReviewSelect, type ReviewScope } from "@/components/ReviewSelect";
import { getReviewWords, getSummary } from "@/lib/progress";
import type { ReviewWord } from "@/lib/types";

// /review            … マイページから。コースと対象を選び、終わったらマイページへ戻る
// /review?course=1   … 単語学習から。コースは固定で対象だけ選び、終わったら学習の続きへ戻る
export default function ReviewSelectPage({ searchParams }: PageProps<"/review">) {
  return (
    <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
      <ReviewSelectData searchParams={searchParams} />
    </Suspense>
  );
}

function countByMastery(words: ReviewWord[]) {
  return {
    weak: words.filter((w) => w.mastery === "weak").length,
    unsure: words.filter((w) => w.mastery === "unsure").length,
  };
}

async function ReviewSelectData({ searchParams }: Pick<PageProps<"/review">, "searchParams">) {
  const courseParam = (await searchParams).course;
  const fixedCourseId = typeof courseParam === "string" ? Number(courseParam) : null;

  // コース名と中断中かどうかは GET /me/summary から
  const summary = await getSummary();
  let courses = summary.courses;
  if (fixedCourseId !== null) {
    courses = courses.filter((c) => c.course_id === fixedCourseId);
    if (courses.length === 0) notFound();
  }
  // 中断中のコースを先に出す（sort は安定なので、それ以外は元の順番のまま）
  courses = [...courses].sort((a, b) => Number(b.status === "in_progress") - Number(a.status === "in_progress"));

  // 選べる範囲ごとに、苦手・微妙の単語数を数えておく（GET /me/review-words）
  const scopes: ReviewScope[] = await Promise.all(
    courses.map(async (c) => ({
      courseId: c.course_id,
      name: c.name,
      interrupted: c.status === "in_progress",
      counts: countByMastery(await getReviewWords({ courseId: c.course_id })),
    })),
  );
  const all = fixedCourseId === null ? countByMastery(await getReviewWords()) : null;

  // /review と /review?course=1 は同じ画面なので、key を変えて選択状態を引き継がないようにする
  return <ReviewSelect key={fixedCourseId ?? "all"} scopes={scopes} allCounts={all} fixedCourseId={fixedCourseId} />;
}
