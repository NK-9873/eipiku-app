// 学習データの取得処理。今はモックデータを読むだけだが、
// API ができたらここを fetch に差し替えれば画面側はそのまま使える。
import { courses, course_words, words, type Course, type Word } from "@/data/mock-data";

export type { Course, Word };

// word_progresses.mastery と同じ値
export type Mastery = "weak" | "unsure" | "mastered";

export function getCourses(): Course[] {
  return courses;
}

export function getCourse(courseId: number): Course | undefined {
  return courses.find((course) => course.id === courseId);
}

export function getCourseWords(courseId: number): Word[] {
  return course_words
    .filter((cw) => cw.course_id === courseId)
    .sort((a, b) => a.display_order - b.display_order)
    .map((cw) => words.find((word) => word.id === cw.word_id))
    .filter((word): word is Word => word !== undefined);
}
