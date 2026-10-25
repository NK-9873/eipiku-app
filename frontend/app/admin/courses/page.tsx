import { Suspense } from "react";
import { CourseList } from "@/components/admin/CourseList";
import { Loading } from "@/components/admin/ui";
import { getCourses, getCourseWords } from "@/lib/learning";

export default function AdminCoursesPage({ searchParams }: PageProps<"/admin/courses">) {
  return (
    <Suspense fallback={<Loading />}>
      <Courses searchParams={searchParams} />
    </Suspense>
  );
}

async function Courses({ searchParams }: Pick<PageProps<"/admin/courses">, "searchParams">) {
  const [{ saved }, courses] = await Promise.all([searchParams, getCourses()]);
  const withWords = await Promise.all(courses.map(async (c) => ({ ...c, words: await getCourseWords(c.id) })));
  return <CourseList courses={withWords} saved={typeof saved === "string" ? saved : undefined} />;
}
