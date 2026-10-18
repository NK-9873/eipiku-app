"use client";

// 管理画面で編集するデータの置き場所。/admin のレイアウトで1回だけ作るので、
// 管理画面の中を移動しても変更は残り、リロードするとモックデータに戻る。
import { createContext, useContext, useState, type ReactNode } from "react";
import {
  courses as initialCourses,
  course_progresses,
  course_words as initialCourseWords,
  users,
  word_progresses,
  words as initialWords,
  type CourseProgress,
  type User,
  type WordProgress,
} from "@/data/mock-data";
import type { Course, CourseWord, Word } from "@/lib/learning";

type AdminStore = {
  courses: Course[];
  words: Word[];
  courseWords: CourseWord[];
  users: User[];
  wordProgresses: WordProgress[];
  courseProgresses: CourseProgress[];
  // id が null なら新規登録。wordIds は出題順に並べたもの
  saveCourse: (course: Omit<Course, "id"> & { id: number | null }, wordIds: number[]) => void;
  deleteCourse: (courseId: number) => void;
  saveWord: (word: Omit<Word, "id"> & { id: number | null }) => void;
  deleteWord: (wordId: number) => void;
};

const AdminStoreContext = createContext<AdminStore | null>(null);

function nextId(rows: { id: number }[]) {
  return Math.max(0, ...rows.map((row) => row.id)) + 1;
}

// 単語を消したあとなどに display_order を 1 から振り直す
function renumber(courseWords: CourseWord[]): CourseWord[] {
  const counters = new Map<number, number>();
  return [...courseWords]
    .sort((a, b) => a.course_id - b.course_id || a.display_order - b.display_order)
    .map((cw) => {
      const order = (counters.get(cw.course_id) ?? 0) + 1;
      counters.set(cw.course_id, order);
      return { ...cw, display_order: order };
    });
}

export function AdminStoreProvider({ children }: { children: ReactNode }) {
  const [courses, setCourses] = useState(initialCourses);
  const [words, setWords] = useState(initialWords);
  const [courseWords, setCourseWords] = useState(initialCourseWords);
  const [courseProgresses, setCourseProgresses] = useState(course_progresses);
  const [wordProgresses, setWordProgresses] = useState(word_progresses);

  const store: AdminStore = {
    courses,
    words,
    courseWords,
    users,
    wordProgresses,
    courseProgresses,

    saveCourse({ id, ...fields }, wordIds) {
      const courseId = id ?? nextId(courses);
      setCourses((prev) =>
        id === null ? [...prev, { id: courseId, ...fields }] : prev.map((c) => (c.id === id ? { id, ...fields } : c)),
      );
      setCourseWords((prev) => [
        ...prev.filter((cw) => cw.course_id !== courseId),
        ...wordIds.map((wordId, i) => ({ course_id: courseId, word_id: wordId, display_order: i + 1 })),
      ]);
    },

    deleteCourse(courseId) {
      setCourses((prev) => prev.filter((c) => c.id !== courseId));
      setCourseWords((prev) => prev.filter((cw) => cw.course_id !== courseId));
      setCourseProgresses((prev) => prev.filter((cp) => cp.course_id !== courseId));
    },

    saveWord({ id, ...fields }) {
      setWords((prev) =>
        id === null ? [...prev, { id: nextId(prev), ...fields }] : prev.map((w) => (w.id === id ? { id, ...fields } : w)),
      );
    },

    deleteWord(wordId) {
      setWords((prev) => prev.filter((w) => w.id !== wordId));
      setCourseWords((prev) => renumber(prev.filter((cw) => cw.word_id !== wordId)));
      setWordProgresses((prev) => prev.filter((wp) => wp.word_id !== wordId));
    },
  };

  return <AdminStoreContext.Provider value={store}>{children}</AdminStoreContext.Provider>;
}

export function useAdminStore() {
  const store = useContext(AdminStoreContext);
  if (!store) throw new Error("useAdminStore は AdminStoreProvider の中で使ってください");
  return store;
}
