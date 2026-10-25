import Link from "next/link";
import type { Course } from "@/lib/types";

type Props = {
  courses: (Course & { total: number; position: number })[];
};

export function CourseList({ courses }: Props) {
  if (courses.length === 0) {
    return <p className="mt-10 text-center text-zinc-500">まだコースがありません</p>;
  }
  return (
    <ul className="mt-6 flex flex-col gap-3">
      {courses.map((course) => (
        <li key={course.id}>
          <Link
            href={`/courses/${course.id}`}
            className="block rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm transition hover:border-orange-300 hover:shadow-md active:scale-[0.99]"
          >
            <div className="flex items-start justify-between gap-3">
              <p className="text-lg font-semibold">{course.name}</p>
              <ProgressBadge position={course.position} total={course.total} />
            </div>
            {course.description && <p className="mt-1 text-sm text-zinc-500">{course.description}</p>}
            <p className="mt-2 text-xs text-zinc-400">{course.total} 語</p>
          </Link>
        </li>
      ))}
    </ul>
  );
}

function ProgressBadge({ position, total }: { position: number; total: number }) {
  if (position <= 0 || total === 0) return null;
  if (position >= total) {
    return <span className="shrink-0 rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-600">完了</span>;
  }
  return (
    <span className="shrink-0 rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-600">
      続きから {position} / {total}
    </span>
  );
}
