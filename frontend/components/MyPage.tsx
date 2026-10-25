import Link from "next/link";
import { MASTERY_LABELS, type Mastery, type Summary } from "@/lib/types";

const STAT_STYLES: Record<Mastery, string> = {
  weak: "text-rose-500",
  unsure: "text-amber-500",
  mastered: "text-emerald-500",
};

const card = "rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm";

// マイページ（GET /me/summary の内容を表示する）
export function MyPage({ userName, summary }: { userName: string; summary: Summary }) {
  // 途中まで学習して止まっているコースのうち、最後に学習したもの
  const resumeCourse = summary.courses
    .filter((c) => c.status === "in_progress")
    .sort((a, b) => (b.updated_at ?? "").localeCompare(a.updated_at ?? ""))[0];
  const reviewCount = summary.weak + summary.unsure;

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <h1 className="text-2xl font-bold">マイページ</h1>
      <p className="mt-1 text-sm text-zinc-500">{userName} さん、今日もがんばりましょう</p>

      {resumeCourse && (
        <Link
          href={`/courses/${resumeCourse.course_id}`}
          className="mt-6 block rounded-2xl bg-orange-500 p-5 text-white shadow-sm transition hover:bg-orange-600 active:scale-[0.99]"
        >
          <p className="text-sm text-orange-100">続きから学習</p>
          <p className="mt-1 text-xl font-bold">{resumeCourse.name}</p>
          <p className="mt-1 text-sm text-orange-100">
            {resumeCourse.current_position} / {resumeCourse.total_words} 語まで学習済み ・ {resumeCourse.current_position + 1}{" "}
            語目から再開 →
          </p>
        </Link>
      )}

      <section className={`${card} mt-4`}>
        <h2 className="font-semibold">学習状況</h2>
        <p className="mt-3">
          <span className="text-3xl font-bold">{summary.learned_words}</span>
          <span className="ml-1 text-zinc-500">/ {summary.total_words} 語 学習済み</span>
        </p>
        <div className="mt-4 grid grid-cols-3 gap-2 text-center">
          {(["mastered", "unsure", "weak"] as const).map((mastery) => (
            <div key={mastery} className="rounded-xl bg-zinc-50 py-3">
              <p className="text-xs text-zinc-500">{MASTERY_LABELS[mastery]}</p>
              <p className={`text-2xl font-bold ${STAT_STYLES[mastery]}`}>{summary[mastery]}</p>
            </div>
          ))}
        </div>
      </section>

      <section className={`${card} mt-4`}>
        <div className="flex items-center justify-between gap-3">
          <div>
            <h2 className="font-semibold">復習</h2>
            <p className="mt-1 text-sm text-zinc-500">
              苦手 {summary.weak} 語 ・ 微妙 {summary.unsure} 語
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
          {summary.courses.map((course) => {
            const done = course.status === "completed";
            return (
              <li key={course.course_id}>
                <div className="flex items-center justify-between gap-3 text-sm">
                  <Link href={`/courses/${course.course_id}`} className="font-medium hover:text-orange-500">
                    {course.name}
                  </Link>
                  <span className="text-zinc-500">{done ? "完了" : `${course.current_position} / ${course.total_words} 語`}</span>
                </div>
                <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-zinc-200">
                  <div
                    className={`h-full ${done ? "bg-emerald-400" : "bg-orange-400"}`}
                    style={{ width: `${course.total_words ? (course.current_position / course.total_words) * 100 : 0}%` }}
                  />
                </div>
              </li>
            );
          })}
        </ul>
      </section>
    </main>
  );
}
