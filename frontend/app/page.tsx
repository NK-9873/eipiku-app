import Link from "next/link";
import { getCourses } from "@/lib/learning";

export default function CourseListPage() {
  const courses = getCourses();

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-10">
      <h1 className="text-3xl font-bold text-orange-500">えいピク</h1>
      <p className="mt-1 text-sm text-zinc-500">絵で覚える英単語帳</p>

      <h2 className="mt-8 mb-3 text-lg font-semibold">コースを選ぶ</h2>
      <ul className="flex flex-col gap-3">
        {courses.map((course) => (
          <li key={course.id}>
            <Link
              href={`/courses/${course.id}`}
              className="block rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm transition hover:border-orange-300 hover:shadow-md active:scale-[0.99]"
            >
              <p className="text-lg font-semibold">{course.name}</p>
              {course.description && (
                <p className="mt-1 text-sm text-zinc-500">{course.description}</p>
              )}
            </Link>
          </li>
        ))}
      </ul>
    </main>
  );
}
