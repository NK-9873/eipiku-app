"use client";

import Link from "next/link";
import type { Course } from "@/lib/learning";
import { useLearner } from "./useLearner";

type Props = {
  courses: (Course & { total: number })[];
};

export function CourseList({ courses }: Props) {
  const learner = useLearner();

  return (
    <ul className="mt-6 flex flex-col gap-3">
      {courses.map((course) => {
        const position = learner?.courseProgresses.find((cp) => cp.course_id === course.id)?.current_position ?? 0;
        return (
          <li key={course.id}>
            <Link
              href={`/courses/${course.id}`}
              className="block rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm transition hover:border-orange-300 hover:shadow-md active:scale-[0.99]"
            >
              <div className="flex items-start justify-between gap-3">
                <p className="text-lg font-semibold">{course.name}</p>
                <ProgressBadge position={position} total={course.total} />
              </div>
              {course.description && <p className="mt-1 text-sm text-zinc-500">{course.description}</p>}
            </Link>
          </li>
        );
      })}
    </ul>
  );
}

function ProgressBadge({ position, total }: { position: number; total: number }) {
  if (position <= 0) return null;
  if (position >= total) {
    return <span className="shrink-0 rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-600">完了</span>;
  }
  return (
    <span className="shrink-0 rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-600">
      続きから {position} / {total}
    </span>
  );
}
