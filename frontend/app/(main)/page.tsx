import { Suspense } from "react";
import { CourseList } from "@/components/CourseList";
import { getCourses, getCourseWords } from "@/lib/learning";

export default function CourseListPage() {
  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <h1 className="text-2xl font-bold">コース一覧</h1>
      <p className="mt-1 text-sm text-zinc-500">学習するコースを選びましょう</p>
      {/* コースは API から毎回取得するので、届くまでは読み込み中を出す */}
      <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
        <Courses />
      </Suspense>
    </main>
  );
}

async function Courses() {
  const courses = await getCourses();
  return <CourseList courses={courses.map((course) => ({ ...course, total: getCourseWords(course.id).length }))} />;
}
