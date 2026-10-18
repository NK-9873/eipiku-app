import { CourseList } from "@/components/CourseList";
import { getCourses, getCourseWords } from "@/lib/learning";

export default function CourseListPage() {
  const courses = getCourses().map((course) => ({ ...course, total: getCourseWords(course.id).length }));

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <h1 className="text-2xl font-bold">コース一覧</h1>
      <p className="mt-1 text-sm text-zinc-500">学習するコースを選びましょう</p>
      <CourseList courses={courses} />
    </main>
  );
}
