import { Suspense } from "react";
import { CourseForm } from "@/components/admin/CourseForm";
import { Loading, NotFoundMessage } from "@/components/admin/ui";
import { getWords, orNotFound } from "@/lib/admin";
import { getCourse, getCourseWords } from "@/lib/learning";

export default function AdminEditCoursePage({ params }: PageProps<"/admin/courses/[courseId]/edit">) {
  return (
    <Suspense fallback={<Loading />}>
      <EditCourse params={params} />
    </Suspense>
  );
}

async function EditCourse({ params }: Pick<PageProps<"/admin/courses/[courseId]/edit">, "params">) {
  const id = Number((await params).courseId);
  const data = await orNotFound(Promise.all([getCourse(id), getCourseWords(id), getWords()]));
  if (!data) return <NotFoundMessage backHref="/admin/courses">コースが見つかりません。削除された可能性があります</NotFoundMessage>;
  const [course, courseWords, allWords] = data;
  return <CourseForm key={course.id} course={course} courseWordIds={courseWords.map((w) => w.id)} allWords={allWords} />;
}
