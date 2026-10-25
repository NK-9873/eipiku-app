import { Suspense } from "react";
import { UserSummary } from "@/components/admin/UserViews";
import { UserForm } from "@/components/admin/UserForm";
import { BackLink, Loading, NotFoundMessage, PageHeader } from "@/components/admin/ui";
import { getUser, getUserProgress, getWords, orNotFound, requireAdmin } from "@/lib/admin";
import { getCourses, getCourseWords } from "@/lib/learning";

export default function AdminUserDetailPage({ params, searchParams }: PageProps<"/admin/users/[userId]">) {
  return (
    <Suspense fallback={<Loading />}>
      <UserDetail params={params} searchParams={searchParams} />
    </Suspense>
  );
}

async function UserDetail({ params, searchParams }: Pick<PageProps<"/admin/users/[userId]">, "params" | "searchParams">) {
  const id = Number((await params).userId);
  const { saved } = await searchParams;
  const [me, user, progress, words, courses] = await Promise.all([
    requireAdmin(),
    orNotFound(getUser(id)),
    orNotFound(getUserProgress(id)),
    getWords(),
    getCourses(),
  ]);
  if (!user || !progress) return <NotFoundMessage backHref="/admin/users">ユーザーが見つかりません</NotFoundMessage>;
  // 進み具合の「何語中」を出すため、コースごとの単語数を数える
  const coursesWithTotal = await Promise.all(courses.map(async (c) => ({ ...c, total: (await getCourseWords(c.id)).length })));

  return (
    <>
      <BackLink href="/admin/users">ユーザー一覧</BackLink>
      <PageHeader title={`${user.name}${user.id === me.id ? "（あなた）" : ""}`} />
      <div className="grid gap-4 md:grid-cols-2">
        <div>
          <UserSummary
            user={user}
            progress={progress}
            words={words}
            courses={coursesWithTotal}
            saved={saved === "created" || saved === "updated" ? saved : undefined}
          />
        </div>
        <div>
          <UserForm key={user.updated_at} user={user} isSelf={user.id === me.id} />
        </div>
      </div>
    </>
  );
}
