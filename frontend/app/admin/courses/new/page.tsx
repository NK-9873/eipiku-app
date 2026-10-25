import { Suspense } from "react";
import { CourseForm } from "@/components/admin/CourseForm";
import { Loading } from "@/components/admin/ui";
import { getWords } from "@/lib/admin";

export default function AdminNewCoursePage() {
  return (
    <Suspense fallback={<Loading />}>
      <NewCourse />
    </Suspense>
  );
}

async function NewCourse() {
  return <CourseForm allWords={await getWords()} />;
}
