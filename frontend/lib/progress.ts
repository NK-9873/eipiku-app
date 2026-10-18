// 学習の記録（word_progresses / course_progresses）の読み書き。
// テーブルと同じ形で localStorage に保存している。API ができたら中身を fetch に差し替える。
import type { CourseProgress, WordProgress } from "@/data/mock-data";
import type { Mastery } from "./learning";
import { now, readStore, writeStore } from "./local-store";

export type { CourseProgress, WordProgress };

const WORD_KEY = "word_progresses";
const COURSE_KEY = "course_progresses";

export async function getWordProgresses(userId: number): Promise<WordProgress[]> {
  return readStore<WordProgress[]>(WORD_KEY, []).filter((wp) => wp.user_id === userId);
}

// UNIQUE (user_id, word_id) なので、あれば上書き・なければ追加する
export async function saveWordProgress(userId: number, wordId: number, mastery: Mastery): Promise<void> {
  const rows = readStore<WordProgress[]>(WORD_KEY, []).filter(
    (wp) => !(wp.user_id === userId && wp.word_id === wordId),
  );
  writeStore(WORD_KEY, [...rows, { user_id: userId, word_id: wordId, mastery, updated_at: now() }]);
}

export async function getCourseProgresses(userId: number): Promise<CourseProgress[]> {
  return readStore<CourseProgress[]>(COURSE_KEY, []).filter((cp) => cp.user_id === userId);
}

// current_position は「何語目まで学習したか」。0 なら未学習、単語数と同じなら最後まで終わっている
export async function saveCourseProgress(userId: number, courseId: number, currentPosition: number): Promise<void> {
  const rows = readStore<CourseProgress[]>(COURSE_KEY, []).filter(
    (cp) => !(cp.user_id === userId && cp.course_id === courseId),
  );
  writeStore(COURSE_KEY, [
    ...rows,
    { user_id: userId, course_id: courseId, current_position: currentPosition, updated_at: now() },
  ]);
}
