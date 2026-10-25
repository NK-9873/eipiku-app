"use client";

import { useEffect, useState } from "react";
import { getCourses, type Course } from "@/lib/learning";
import { getCourseProgresses, getWordProgresses, type CourseProgress, type WordProgress } from "@/lib/progress";
import { CURRENT_USER_ID, getCurrentUser, type User } from "@/lib/user";

export type LearnerData = {
  user: User;
  courses: Course[];
  wordProgresses: WordProgress[];
  courseProgresses: CourseProgress[];
};

// ログイン中のユーザー・コース一覧・学習の記録を読み込む。読み込みが終わるまでは null。
// 読み込みに失敗したら、描画中に投げ直して error.tsx に表示させる
export function useLearner(): LearnerData | null {
  const [data, setData] = useState<LearnerData | null>(null);
  const [error, setError] = useState<unknown>(null);

  useEffect(() => {
    let cancelled = false;
    Promise.all([
      getCurrentUser(),
      getCourses(),
      getWordProgresses(CURRENT_USER_ID),
      getCourseProgresses(CURRENT_USER_ID),
    ]).then(
      ([user, courses, wordProgresses, courseProgresses]) => {
        if (!cancelled) setData({ user, courses, wordProgresses, courseProgresses });
      },
      (e) => {
        if (!cancelled) setError(e);
      },
    );
    return () => {
      cancelled = true;
    };
  }, []);

  if (error) throw error;
  return data;
}
