// コース・単語の取得（learning-service）。サーバー側だけで使う
import "server-only";
import { apiFetch } from "./api";
import type { Course, CourseWord } from "./types";

export function getCourses() {
  return apiFetch<Course[]>("learning", "/courses");
}

export function getCourse(courseId: number) {
  return apiFetch<Course>("learning", `/courses/${courseId}`);
}

// コースの単語を出題順に返す
export function getCourseWords(courseId: number) {
  return apiFetch<CourseWord[]>("learning", `/courses/${courseId}/words`);
}
