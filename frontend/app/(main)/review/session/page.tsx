import { Suspense } from "react";
import { ReviewSession } from "@/components/ReviewSession";
import { getReviewWords } from "@/lib/progress";
import { REVIEW_TARGETS } from "@/lib/types";
import { getMe } from "@/lib/user";

// URL: /review/session?course=all|<id>&target=weak|unsure|both&back=mypage|course
export default function ReviewSessionPage({ searchParams }: PageProps<"/review/session">) {
  return (
    <Suspense fallback={<p className="py-10 text-center text-sm text-zinc-400">読み込み中…</p>}>
      <ReviewSessionData searchParams={searchParams} />
    </Suspense>
  );
}

async function ReviewSessionData({ searchParams }: Pick<PageProps<"/review/session">, "searchParams">) {
  const params = await searchParams;
  const courseId = typeof params.course === "string" && params.course !== "all" ? Number(params.course) : undefined;
  const target = REVIEW_TARGETS.find((t) => t.value === params.target) ?? REVIEW_TARGETS[2];
  // 単語学習から来たときは学習の続き（中断位置）へ、マイページから来たときはマイページへ戻る
  const backHref = params.back === "course" && courseId ? `/courses/${courseId}` : "/mypage";

  // 出題する単語は開始時点で決める（評価して対象から外れても、今回の復習では最後まで出す）
  const [words, me] = await Promise.all([getReviewWords({ courseId, masteries: target.masteries }), getMe()]);

  return (
    <ReviewSession
      // 条件が変わったら最初からやり直す（同じ画面のまま別の条件の復習を開いたとき）
      key={`${courseId ?? "all"}-${target.value}`}
      words={words}
      targetLabel={target.label}
      imageType={me.image_type}
      backHref={backHref}
      backLabel={backHref === "/mypage" ? "マイページに戻る" : "学習を再開する"}
    />
  );
}
