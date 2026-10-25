import Link from "next/link";
import { IMAGE_TYPE_LABELS, MASTERY_LABELS, type Course, type Mastery, type User, type UserProgress, type Word } from "@/lib/types";
import { cardClass, PageHeader, primaryButtonClass, SavedNotice } from "./ui";

const MASTERY_COLORS: Record<Mastery, string> = {
  weak: "bg-rose-50 text-rose-600",
  unsure: "bg-amber-50 text-amber-600",
  mastered: "bg-emerald-50 text-emerald-600",
};

export function RoleBadge({ role }: { role: User["role"] }) {
  return role === "admin" ? (
    <span className="rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-600">管理者</span>
  ) : (
    <span className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs text-zinc-600">一般</span>
  );
}

// ユーザー一覧（GET /users）
export function UserList({ users, meId }: { users: User[]; meId: number }) {
  return (
    <>
      <PageHeader
        title="ユーザー管理"
        action={
          <Link href="/admin/users/new" className={primaryButtonClass}>
            ＋ ユーザーを作成
          </Link>
        }
      />
      <ul className="flex flex-col gap-3">
        {users.map((user) => (
          <li key={user.id}>
            <Link
              href={`/admin/users/${user.id}`}
              className={`${cardClass} flex items-center justify-between gap-3 p-5 transition hover:border-orange-300 hover:shadow-md`}
            >
              <div className="min-w-0">
                <p className="flex flex-wrap items-center gap-2">
                  <span className="text-lg font-semibold">{user.name}</span>
                  <RoleBadge role={user.role} />
                  {user.id === meId && <span className="text-xs text-zinc-400">（あなた）</span>}
                </p>
                <p className="mt-0.5 truncate text-sm text-zinc-500">{user.email}</p>
              </div>
              <span className="shrink-0 text-sm text-zinc-400">詳細・編集 ›</span>
            </Link>
          </li>
        ))}
      </ul>
    </>
  );
}

type CourseWithTotal = Course & { total: number };

// ユーザーの詳細（基本情報と、GET /admin/users/{id}/progress の学習の記録）
export function UserSummary({
  user,
  progress,
  words,
  courses,
  saved,
}: {
  user: User;
  progress: UserProgress;
  words: Word[];
  courses: CourseWithTotal[];
  saved?: "created" | "updated";
}) {
  return (
    <>
      {saved && <SavedNotice>{saved === "created" ? "ユーザーを作成しました" : "変更を保存しました"}</SavedNotice>}

      <section className={`${cardClass} p-5`}>
        <h2 className="font-semibold">基本情報</h2>
        <dl className="mt-3 grid grid-cols-[7rem_1fr] gap-y-2 text-sm">
          <dt className="text-zinc-500">権限</dt>
          <dd>
            <RoleBadge role={user.role} />
          </dd>
          <dt className="text-zinc-500">画像の種類</dt>
          <dd>{IMAGE_TYPE_LABELS[user.image_type]}</dd>
          <dt className="text-zinc-500">登録日時</dt>
          <dd>{new Date(user.created_at).toLocaleString("ja-JP")}</dd>
        </dl>
      </section>

      <section className={`${cardClass} mt-4 p-5`}>
        <h2 className="font-semibold">コースの進み具合</h2>
        {progress.course_progresses.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-400">まだ学習していません</p>
        ) : (
          <ul className="mt-3 flex flex-col gap-3">
            {progress.course_progresses.map((cp) => {
              const course = courses.find((c) => c.id === cp.course_id);
              const total = course?.total ?? 0;
              const position = Math.min(cp.current_position, total);
              return (
                <li key={cp.course_id}>
                  <div className="flex justify-between text-sm">
                    <span>{course?.name ?? "（削除されたコース）"}</span>
                    <span className="text-zinc-500">{total > 0 && position >= total ? "完了" : `${position} / ${total} 語`}</span>
                  </div>
                  <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-zinc-200">
                    <div className="h-full bg-orange-400" style={{ width: `${total ? (position / total) * 100 : 0}%` }} />
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </section>

      <section className={`${cardClass} mt-4 p-5`}>
        <h2 className="font-semibold">
          単語の覚え具合 <span className="text-sm font-normal text-zinc-500">（{progress.word_progresses.length} 語）</span>
        </h2>
        {progress.word_progresses.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-400">まだ評価した単語はありません</p>
        ) : (
          <ul className="mt-3 flex flex-col divide-y divide-zinc-100">
            {progress.word_progresses.map((wp) => {
              const word = words.find((w) => w.id === wp.word_id);
              return (
                <li key={wp.word_id} className="flex items-center justify-between gap-3 py-2 text-sm">
                  <span>
                    {word?.english ?? "（削除された単語）"} <span className="text-zinc-500">{word?.japanese}</span>
                  </span>
                  <span className={`rounded-full px-2 py-0.5 text-xs font-semibold ${MASTERY_COLORS[wp.mastery]}`}>
                    {MASTERY_LABELS[wp.mastery]}
                  </span>
                </li>
              );
            })}
          </ul>
        )}
      </section>
    </>
  );
}
