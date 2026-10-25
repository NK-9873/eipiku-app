"use client";

import Link from "next/link";
import { pickCourseWords, useAdminStore } from "./AdminStore";
import { cardClass, dangerButtonClass, PageHeader, primaryButtonClass, smallButtonClass } from "./ui";

export function CourseList() {
  const { courses, courseWords, words, deleteCourse } = useAdminStore();

  function handleDelete(courseId: number, name: string) {
    if (window.confirm(`コース「${name}」を削除しますか？`)) deleteCourse(courseId);
  }

  return (
    <>
      <PageHeader
        title="コース管理"
        action={
          <Link href="/admin/courses/new" className={primaryButtonClass}>
            ＋ 新規登録
          </Link>
        }
      />

      {courses.length === 0 ? (
        <p className={`${cardClass} p-8 text-center text-zinc-500`}>コースがありません</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {courses.map((course) => {
            const courseWordList = pickCourseWords(course.id, courseWords, words);
            return (
              <li key={course.id} className={`${cardClass} p-5`}>
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="text-lg font-semibold">{course.name}</p>
                    {course.description && <p className="mt-1 text-sm text-zinc-500">{course.description}</p>}
                  </div>
                  <div className="flex shrink-0 gap-1">
                    <Link href={`/admin/courses/${course.id}/edit`} className={smallButtonClass}>
                      編集
                    </Link>
                    <button type="button" onClick={() => handleDelete(course.id, course.name)} className={dangerButtonClass}>
                      削除
                    </button>
                  </div>
                </div>
                <p className="mt-3 text-sm text-zinc-600">
                  <span className="mr-2 rounded-full bg-orange-50 px-2 py-0.5 text-xs font-semibold text-orange-600">
                    {courseWordList.length} 語
                  </span>
                  {courseWordList.map((w) => w.english).join(" → ") || "単語が登録されていません"}
                </p>
              </li>
            );
          })}
        </ul>
      )}
    </>
  );
}
