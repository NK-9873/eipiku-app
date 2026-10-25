import { Suspense } from "react";
import { CourseList } from "@/components/CourseList";
import { getCourses } from "@/lib/learning";
import { getSummary } from "@/lib/progress";

export default function CourseListPage() {
  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <h1 className="text-2xl font-bold">コース一覧</h1>
      <p className="mt-1 text-sm text-zinc-500">学習するコースを選びましょう</p>
      {/* コースと進み具合は API から毎回取得するので、届くまでは読み込み中を出す */}
      <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
        <Courses />
      </Suspense>
    </main>
  );
}

async function Courses() {
  // コースの説明は GET /courses、進み具合（単語数・中断位置）は GET /me/summary から取る
  const [courses, summary] = await Promise.all([getCourses(), getSummary()]);
  return (
    <CourseList
      courses={courses.map((course) => {
        const progress = summary.courses.find((c) => c.course_id === course.id);
        return { ...course, total: progress?.total_words ?? 0, position: progress?.current_position ?? 0 };
      })}
    />
  );
}
