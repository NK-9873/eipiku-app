// 学習の記録（習得度・進み具合）の取得（learning-service）。サーバー側だけで使う。
// 保存は画面から呼ぶので lib/actions.ts にある
import "server-only";
import { apiFetch } from "./api";
import type { CourseProgress, Mastery, ReviewWord, Summary } from "./types";

export function getSummary() {
  return apiFetch<Summary>("learning", "/me/summary");
}

export function getCourseProgresses() {
  return apiFetch<CourseProgress[]>("learning", "/me/course-progresses");
}

// 復習対象の単語。courseId を省略すると全コース、masteries を省略すると苦手と微妙
export function getReviewWords(filter: { courseId?: number; masteries?: Mastery[] } = {}) {
  const params = new URLSearchParams();
  if (filter.courseId) params.set("course_id", String(filter.courseId));
  if (filter.masteries?.length) params.set("mastery", filter.masteries.join(","));
  const query = params.size ? `?${params}` : "";
  return apiFetch<ReviewWord[]>("learning", `/me/review-words${query}`);
}
