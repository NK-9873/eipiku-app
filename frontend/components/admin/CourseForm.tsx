"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";
import type { Course } from "@/data/mock-data";
import { pickCourseWords, useAdminStore } from "./AdminStore";
import {
  BackLink,
  cardClass,
  dangerButtonClass,
  Field,
  inputClass,
  NotFoundMessage,
  PageHeader,
  primaryButtonClass,
  secondaryButtonClass,
  smallButtonClass,
} from "./ui";

// /admin/courses/[courseId]/edit 用。URL の id からコースを探してフォームに渡す
export function EditCourse() {
  const { courseId } = useParams<{ courseId: string }>();
  const { courses } = useAdminStore();
  const course = courses.find((c) => c.id === Number(courseId));

  if (!course) {
    return <NotFoundMessage backHref="/admin/courses">コースが見つかりません</NotFoundMessage>;
  }
  return <CourseForm course={course} />;
}

export function CourseForm({ course }: { course?: Course }) {
  const router = useRouter();
  const { words, courseWords, saveCourse } = useAdminStore();

  const [name, setName] = useState(course?.name ?? "");
  const [description, setDescription] = useState(course?.description ?? "");
  // 出題順に並んだ、コースに含める単語の id
  const [selectedIds, setSelectedIds] = useState<number[]>(() =>
    course ? pickCourseWords(course.id, courseWords, words).map((w) => w.id) : [],
  );
  const [error, setError] = useState("");

  const selectedWords = selectedIds
    .map((id) => words.find((w) => w.id === id))
    .filter((w) => w !== undefined);
  const unselectedWords = words.filter((w) => !selectedIds.includes(w.id));

  function move(index: number, offset: -1 | 1) {
    setSelectedIds((prev) => {
      const next = [...prev];
      [next[index], next[index + offset]] = [next[index + offset], next[index]];
      return next;
    });
  }

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!name.trim()) {
      setError("コース名を入力してください");
      return;
    }
    saveCourse(
      { id: course?.id ?? null, name: name.trim(), description: description.trim() || null },
      selectedIds,
    );
    router.push("/admin/courses");
  }

  return (
    <form onSubmit={handleSubmit}>
      <BackLink href="/admin/courses">コース一覧</BackLink>
      <PageHeader title={course ? "コースを編集" : "コースを登録"} />

      <div className={`${cardClass} flex flex-col gap-4 p-5`}>
        <Field label="コース名" required>
          <input value={name} onChange={(e) => setName(e.target.value)} maxLength={100} className={inputClass} />
        </Field>
        <Field label="説明">
          <input
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            maxLength={255}
            className={inputClass}
          />
        </Field>
      </div>

      <div className="mt-5 grid gap-4 sm:grid-cols-2">
        <section className={`${cardClass} p-5`}>
          <h2 className="font-semibold">
            含める単語 <span className="text-sm font-normal text-zinc-500">（上から出題）</span>
          </h2>
          {selectedWords.length === 0 ? (
            <p className="mt-3 text-sm text-zinc-400">右の一覧から単語を追加してください</p>
          ) : (
            <ol className="mt-3 flex flex-col gap-2">
              {selectedWords.map((word, i) => (
                <li key={word.id} className="flex items-center gap-2 rounded-xl bg-orange-50 px-3 py-2">
                  <span className="w-5 text-sm font-semibold text-orange-500">{i + 1}</span>
                  <span className="min-w-0 flex-1 truncate">
                    {word.english} <span className="text-sm text-zinc-500">{word.japanese}</span>
                  </span>
                  <button type="button" onClick={() => move(i, -1)} disabled={i === 0} className={smallButtonClass} aria-label="上へ">
                    ↑
                  </button>
                  <button
                    type="button"
                    onClick={() => move(i, 1)}
                    disabled={i === selectedWords.length - 1}
                    className={smallButtonClass}
                    aria-label="下へ"
                  >
                    ↓
                  </button>
                  <button
                    type="button"
                    onClick={() => setSelectedIds((prev) => prev.filter((id) => id !== word.id))}
                    className={dangerButtonClass}
                  >
                    外す
                  </button>
                </li>
              ))}
            </ol>
          )}
        </section>

        <section className={`${cardClass} p-5`}>
          <h2 className="font-semibold">追加できる単語</h2>
          {unselectedWords.length === 0 ? (
            <p className="mt-3 text-sm text-zinc-400">すべての単語が含まれています</p>
          ) : (
            <ul className="mt-3 flex flex-col gap-2">
              {unselectedWords.map((word) => (
                <li key={word.id} className="flex items-center gap-2 rounded-xl border border-zinc-100 px-3 py-2">
                  <span className="min-w-0 flex-1 truncate">
                    {word.english} <span className="text-sm text-zinc-500">{word.japanese}</span>
                  </span>
                  <button
                    type="button"
                    onClick={() => setSelectedIds((prev) => [...prev, word.id])}
                    className={smallButtonClass}
                  >
                    追加
                  </button>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      {error && <p className="mt-4 text-sm text-rose-600">{error}</p>}
      <div className="mt-6 flex justify-end gap-2">
        <Link href="/admin/courses" className={secondaryButtonClass}>
          キャンセル
        </Link>
        <button type="submit" className={primaryButtonClass}>
          {course ? "保存する" : "登録する"}
        </button>
      </div>
    </form>
  );
}
