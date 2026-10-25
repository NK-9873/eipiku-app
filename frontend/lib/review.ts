// 復習の対象を決める処理
import { getCourseWords, type Mastery, type Word } from "./learning";
import type { CourseProgress, WordProgress } from "./progress";

export type ReviewTarget = "weak" | "unsure" | "both";

export const REVIEW_TARGETS: { value: ReviewTarget; label: string }[] = [
  { value: "weak", label: "苦手" },
  { value: "unsure", label: "微妙" },
  { value: "both", label: "両方" },
];

function matches(target: ReviewTarget, mastery: Mastery) {
  return target === "both" ? mastery !== "mastered" : mastery === target;
}

// 指定したコースの単語から選ぶ。同じ単語が複数のコースにあっても1回だけ出す
export function pickReviewWords(courseIds: number[], target: ReviewTarget, wordProgresses: WordProgress[]): Word[] {
  const courseWords = courseIds.flatMap((id) => getCourseWords(id));
  const unique = [...new Map(courseWords.map((w) => [w.id, w])).values()];
  return unique.filter((word) => {
    const progress = wordProgresses.find((wp) => wp.word_id === word.id);
    return progress !== undefined && matches(target, progress.mastery);
  });
}

// 途中まで学習して止まっているか（最初から・最後まで終わったものは含めない）
export function isInterrupted(courseId: number, courseProgresses: CourseProgress[]) {
  const position = courseProgresses.find((cp) => cp.course_id === courseId)?.current_position ?? 0;
  return position > 0 && position < getCourseWords(courseId).length;
}
