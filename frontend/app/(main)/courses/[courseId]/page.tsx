import { notFound } from "next/navigation";
import { Suspense } from "react";
import { WordStudy } from "@/components/WordStudy";
import { ApiError } from "@/lib/api";
import { getCourse, getCourseWords } from "@/lib/learning";
import { getCourseProgresses } from "@/lib/progress";
import { getMe } from "@/lib/user";

export default function WordStudyPage({ params }: PageProps<"/courses/[courseId]">) {
  // params とログイン中のユーザーの記録は Suspense の中で読む（cacheComponents で画面遷移をすぐ始めるため）
  return (
    <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
      <CourseWordStudy params={params} />
    </Suspense>
  );
}

async function CourseWordStudy({ params }: Pick<PageProps<"/courses/[courseId]">, "params">) {
  const courseId = Number((await params).courseId);
  if (!Number.isInteger(courseId) || courseId <= 0) notFound();

  const [course, words, progresses, me] = await Promise.all([
    getCourse(courseId).catch((e) => {
      if (e instanceof ApiError && e.status === 404) notFound();
      throw e;
    }),
    getCourseWords(courseId),
    getCourseProgresses(),
    getMe(),
  ]);

  // 途中まで学習して止まっていれば続きから。未学習・最後まで終わったコースは最初から
  const saved = progresses.find((p) => p.course_id === courseId)?.current_position ?? 0;
  const startIndex = saved > 0 && saved < words.length ? saved : 0;

  // 別のコースに移ったら最初から作り直す（同じ画面のまま進み具合を引き継がないように）
  return <WordStudy key={course.id} course={course} words={words} startIndex={startIndex} imageType={me.image_type} />;
}
