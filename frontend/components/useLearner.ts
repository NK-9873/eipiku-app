"use client";

import { useEffect, useState } from "react";
import { getCourseProgresses, getWordProgresses, type CourseProgress, type WordProgress } from "@/lib/progress";
import { CURRENT_USER_ID, getCurrentUser, type User } from "@/lib/user";

export type LearnerData = {
  user: User;
  wordProgresses: WordProgress[];
  courseProgresses: CourseProgress[];
};

// ログイン中のユーザーと学習の記録を読み込む。読み込みが終わるまでは null
export function useLearner(): LearnerData | null {
  const [data, setData] = useState<LearnerData | null>(null);

  useEffect(() => {
    let cancelled = false;
    Promise.all([
      getCurrentUser(),
      getWordProgresses(CURRENT_USER_ID),
      getCourseProgresses(CURRENT_USER_ID),
    ]).then(([user, wordProgresses, courseProgresses]) => {
      if (!cancelled) setData({ user, wordProgresses, courseProgresses });
    });
    return () => {
      cancelled = true;
    };
  }, []);

  return data;
}
