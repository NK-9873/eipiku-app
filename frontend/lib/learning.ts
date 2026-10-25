// 学習データの取得処理。コース一覧は learning-service の API から、
// それ以外はまだモックデータから読んでいる。API ができたものから fetch に差し替える。
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

const API_URL = process.env.NEXT_PUBLIC_LEARNING_API_URL;

// GET /courses（learning-service）
export async function getCourses(): Promise<Course[]> {
  if (!API_URL) {
    throw new Error("NEXT_PUBLIC_LEARNING_API_URL が設定されていません（frontend/.env.local を確認してください）");
  }
  const res = await fetch(`${API_URL}/courses`);
  if (!res.ok) {
    throw new Error(`コース一覧の取得に失敗しました（HTTP ${res.status}）`);
  }
  return res.json();
}

export function getWords(): Word[] {
  return words;
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
