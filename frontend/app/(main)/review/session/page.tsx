import { Suspense } from "react";
import { ReviewSession } from "@/components/ReviewSession";

export default function ReviewSessionPage() {
  return (
    <Suspense>
      <ReviewSession />
    </Suspense>
  );
}
