import { Suspense } from "react";
import { EditCourse } from "@/components/admin/CourseForm";
import { Loading } from "@/components/admin/ui";

export default function AdminEditCoursePage() {
  // 管理画面で追加したコースの id は事前に分からないので、useParams を Suspense で包む
  return (
    <Suspense fallback={<Loading />}>
      <EditCourse />
    </Suspense>
  );
}
