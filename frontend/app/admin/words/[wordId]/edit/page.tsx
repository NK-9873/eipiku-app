import { Suspense } from "react";
import { EditWord } from "@/components/admin/WordForm";
import { Loading } from "@/components/admin/ui";

export default function AdminEditWordPage() {
  // 管理画面で追加した単語の id は事前に分からないので、useParams を Suspense で包む
  return (
    <Suspense fallback={<Loading />}>
      <EditWord />
    </Suspense>
  );
}
