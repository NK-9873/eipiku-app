import Link from "next/link";
import { deleteCourse } from "@/lib/admin-actions";
import type { Course, CourseWord } from "@/lib/types";
import { DeleteButton } from "./DeleteButton";
import { Thumbnail } from "./Thumbnail";
import { cardClass, PageHeader, primaryButtonClass, SavedNotice, smallButtonClass } from "./ui";

// コース一覧（GET /courses と、コースごとの GET /courses/{id}/words）
export function CourseList({ courses, saved }: { courses: (Course & { words: CourseWord[] })[]; saved?: string }) {
  return (
    <>
      <PageHeader
        title="コース管理"
        action={
          <Link href="/admin/courses/new" className={primaryButtonClass}>
            ＋ コースを登録
          </Link>
        }
      />
      {saved && <SavedNotice>コース「{saved}」を保存しました</SavedNotice>}

      {courses.length === 0 ? (
        <p className={`${cardClass} p-8 text-center text-zinc-500`}>まだコースがありません。「＋ コースを登録」から追加してください</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {courses.map((course) => (
            <li key={course.id} className={`${cardClass} p-5`}>
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-lg font-semibold">{course.name}</p>
                  {course.description && <p className="mt-1 text-sm text-zinc-500">{course.description}</p>}
                </div>
                <div className="flex shrink-0 items-start gap-1">
                  <Link href={`/admin/courses/${course.id}/edit`} className={smallButtonClass}>
                    編集
                  </Link>
                  <DeleteButton
                    confirmMessage={`コース「${course.name}」を削除しますか？\n単語そのものは消えませんが、このコースの学習の進み具合は消えます。元に戻せません。`}
                    action={deleteCourse.bind(null, course.id)}
                  />
                </div>
              </div>
              <p className="mt-3 text-sm text-zinc-600">
                <span className="mr-2 rounded-full bg-orange-50 px-2 py-0.5 text-xs font-semibold text-orange-600">
                  {course.words.length} 語
                </span>
                {course.words.length === 0 ? "単語がまだ入っていません" : "出題順："}
              </p>
              {course.words.length > 0 && (
                <ol className="mt-2 flex flex-wrap gap-2">
                  {course.words.map((w) => (
                    <li key={w.id} className="flex items-center gap-1.5 rounded-lg bg-zinc-50 py-1 pl-1 pr-2 text-sm">
                      <Thumbnail src={w.photo_path ?? w.illustration_path} alt={w.english} size="h-7 w-7" />
                      {w.display_order}. {w.english}
                    </li>
                  ))}
                </ol>
              )}
            </li>
          ))}
        </ul>
      )}
    </>
  );
}
