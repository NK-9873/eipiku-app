// 学習データの取得処理。今はモックデータを読むだけだが、
// API ができたらここを fetch に差し替えれば画面側はそのまま使える。
import {
  courses,
  course_words,
  words,
  type Course,
  type CourseWord,
  type Word,
} from "@/data/mock-data";

export type { Course, CourseWord, Word };

// word_progresses.mastery と同じ値
export type Mastery = "weak" | "unsure" | "mastered";

export const MASTERY_LABELS: Record<Mastery, string> = {
  weak: "苦手",
  unsure: "微妙",
  mastered: "覚えた",
};

export function getCourses(): Course[] {
  return courses;
}

export function getCourse(courseId: number): Course | undefined {
  return courses.find((course) => course.id === courseId);
}

export function getCourseWords(courseId: number): Word[] {
  return pickCourseWords(courseId, course_words, words);
}

// course_words を出題順に並べて単語に変換する。管理画面では編集中のデータを渡して使う。
export function pickCourseWords(courseId: number, courseWords: CourseWord[], allWords: Word[]): Word[] {
  return courseWords
    .filter((cw) => cw.course_id === courseId)
    .sort((a, b) => a.display_order - b.display_order)
    .map((cw) => allWords.find((word) => word.id === cw.word_id))
    .filter((word): word is Word => word !== undefined);
}
