"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import type { User } from "@/data/mock-data";
import { MASTERY_LABELS, pickCourseWords, type Mastery } from "@/lib/learning";
import { useAdminStore } from "./AdminStore";
import { BackLink, cardClass, NotFoundMessage, PageHeader } from "./ui";

const IMAGE_TYPE_LABELS: Record<User["image_type"], string> = {
  photo: "写真",
  illustration: "イラスト",
};

const MASTERY_COLORS: Record<Mastery, string> = {
  weak: "bg-rose-50 text-rose-600",
  unsure: "bg-amber-50 text-amber-600",
  mastered: "bg-emerald-50 text-emerald-600",
};

function RoleBadge({ role }: { role: User["role"] }) {
  return role === "admin" ? (
    <span className="rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-600">管理者</span>
  ) : (
    <span className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs text-zinc-600">一般</span>
  );
}

export function UserList() {
  const { users } = useAdminStore();

  return (
    <>
      <PageHeader title="ユーザー管理" />
      <ul className="flex flex-col gap-3">
        {users.map((user) => (
          <li key={user.id}>
            <Link
              href={`/admin/users/${user.id}`}
              className={`${cardClass} flex items-center justify-between gap-3 p-5 transition hover:border-orange-300 hover:shadow-md`}
            >
              <div className="min-w-0">
                <p className="flex items-center gap-2">
                  <span className="text-lg font-semibold">{user.name}</span>
                  <RoleBadge role={user.role} />
                </p>
                <p className="mt-0.5 truncate text-sm text-zinc-500">{user.email}</p>
              </div>
              <span className="shrink-0 text-zinc-300">›</span>
            </Link>
          </li>
        ))}
      </ul>
    </>
  );
}

export function UserDetail() {
  const { userId } = useParams<{ userId: string }>();
  const { users, courses, courseWords, words, courseProgresses, wordProgresses } = useAdminStore();
  const user = users.find((u) => u.id === Number(userId));

  if (!user) {
    return <NotFoundMessage backHref="/admin/users">ユーザーが見つかりません</NotFoundMessage>;
  }

  const progresses = courseProgresses.filter((cp) => cp.user_id === user.id);
  const masteries = wordProgresses.filter((wp) => wp.user_id === user.id);

  return (
    <>
      <BackLink href="/admin/users">ユーザー一覧</BackLink>
      <PageHeader title={user.name} />

      <section className={`${cardClass} p-5`}>
        <h2 className="font-semibold">基本情報</h2>
        <dl className="mt-3 grid grid-cols-[7rem_1fr] gap-y-2 text-sm">
          <dt className="text-zinc-500">ID</dt>
          <dd>{user.id}</dd>
          <dt className="text-zinc-500">メールアドレス</dt>
          <dd className="break-all">{user.email}</dd>
          <dt className="text-zinc-500">権限</dt>
          <dd>
            <RoleBadge role={user.role} />
          </dd>
          <dt className="text-zinc-500">画像の種類</dt>
          <dd>{IMAGE_TYPE_LABELS[user.image_type]}</dd>
          <dt className="text-zinc-500">登録日時</dt>
          <dd>{user.created_at}</dd>
        </dl>
      </section>

      <section className={`${cardClass} mt-4 p-5`}>
        <h2 className="font-semibold">コースの進み具合</h2>
        {progresses.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-400">まだ学習していません</p>
        ) : (
          <ul className="mt-3 flex flex-col gap-3">
            {progresses.map((cp) => {
              const course = courses.find((c) => c.id === cp.course_id);
              const total = pickCourseWords(cp.course_id, courseWords, words).length;
              const position = Math.min(cp.current_position, total);
              return (
                <li key={cp.course_id}>
                  <div className="flex justify-between text-sm">
                    <span>{course?.name ?? "（削除されたコース）"}</span>
                    <span className="text-zinc-500">
                      {position} / {total} 語
                    </span>
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
        <h2 className="font-semibold">単語の覚え具合</h2>
        {masteries.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-400">まだ評価した単語はありません</p>
        ) : (
          <ul className="mt-3 flex flex-col divide-y divide-zinc-100">
            {masteries.map((wp) => {
              const word = words.find((w) => w.id === wp.word_id);
              return (
                <li key={wp.word_id} className="flex items-center justify-between gap-3 py-2 text-sm">
                  <span>
                    {word?.english} <span className="text-zinc-500">{word?.japanese}</span>
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
