import { notFound } from "next/navigation";
import { Suspense } from "react";
import { getCourse, getCourseWords } from "@/lib/learning";
import { WordStudy } from "@/components/WordStudy";

export default function WordStudyPage({ params }: PageProps<"/courses/[courseId]">) {
  // params は Suspense の中で読む（cacheComponents で画面遷移をすぐ始めるため）
  return (
    <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
      <CourseWordStudy params={params} />
    </Suspense>
  );
}

async function CourseWordStudy({ params }: Pick<PageProps<"/courses/[courseId]">, "params">) {
  const { courseId } = await params;
  const course = getCourse(Number(courseId));
  if (!course) notFound();

  return <WordStudy course={course} words={getCourseWords(course.id)} />;
}
