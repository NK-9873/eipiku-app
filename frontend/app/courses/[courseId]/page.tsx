import { notFound } from "next/navigation";
import { getCourse, getCourses, getCourseWords } from "@/lib/learning";
import { WordStudy } from "@/components/WordStudy";

export function generateStaticParams() {
  return getCourses().map((course) => ({ courseId: String(course.id) }));
}

export default async function WordStudyPage({ params }: PageProps<"/courses/[courseId]">) {
  const { courseId } = await params;
  const course = getCourse(Number(courseId));
  if (!course) notFound();

  return <WordStudy course={course} words={getCourseWords(course.id)} />;
}
