"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useTransition, type FormEvent } from "react";
import { saveCourse } from "@/lib/admin-actions";
import type { Course, FormState, Word } from "@/lib/types";
import {
  BackLink,
  cardClass,
  dangerButtonClass,
  ErrorBox,
  Field,
  inputClass,
  PageHeader,
  primaryButtonClass,
  secondaryButtonClass,
  smallButtonClass,
} from "./ui";
import { Thumbnail } from "./Thumbnail";

type Props = {
  course?: Course;
  courseWordIds?: number[]; // 出題順に並んだ、コースに入っている単語
  allWords: Word[];
};

// コースの登録・編集。含める単語と出題順も一緒に保存する
export function CourseForm({ course, courseWordIds = [], allWords }: Props) {
  const router = useRouter();
  const [name, setName] = useState(course?.name ?? "");
  const [description, setDescription] = useState(course?.description ?? "");
  const [selectedIds, setSelectedIds] = useState<number[]>(courseWordIds);
  const [filter, setFilter] = useState("");
  const [state, setState] = useState<FormState>({});
  const [pending, startTransition] = useTransition();

  const selectedWords = selectedIds.map((id) => allWords.find((w) => w.id === id)).filter((w) => w !== undefined);
  const keyword = filter.trim().toLowerCase();
  const unselectedWords = allWords.filter(
    (w) => !selectedIds.includes(w.id) && (!keyword || w.english.toLowerCase().includes(keyword) || w.japanese.includes(keyword)),
  );

  function move(index: number, offset: -1 | 1) {
    setSelectedIds((prev) => {
      const next = [...prev];
      [next[index], next[index + offset]] = [next[index + offset], next[index]];
      return next;
    });
  }

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    startTransition(async () => {
      const result = await saveCourse(course?.id ?? null, { name, description, wordIds: selectedIds });
      setState(result);
      // 新しいコースは登録できたが単語の保存に失敗したとき、編集画面に移ってやり直せるようにする
      if (!course && result.courseId) router.replace(`/admin/courses/${result.courseId}/edit`);
    });
  }

  return (
    <form onSubmit={handleSubmit}>
      <BackLink href="/admin/courses">コース一覧</BackLink>
      <PageHeader title={course ? `コースを編集：${course.name}` : "コースを登録"} />

      <div className={`${cardClass} flex flex-col gap-4 p-5`}>
        <Field label="コース名" required>
          <input value={name} onChange={(e) => setName(e.target.value)} required maxLength={100} placeholder="例：動物" className={inputClass} />
        </Field>
        <Field label="説明">
          <input
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            maxLength={255}
            placeholder="例：身近な動物の名前を覚えよう"
            className={inputClass}
          />
        </Field>
      </div>

      <div className="mt-5 grid gap-4 sm:grid-cols-2">
        <section className={`${cardClass} p-5`}>
          <h2 className="font-semibold">
            コースに入れる単語 <span className="text-sm font-normal text-zinc-500">（上から順に出題・{selectedWords.length} 語）</span>
          </h2>
          {selectedWords.length === 0 ? (
            <p className="mt-3 rounded-xl bg-zinc-50 p-4 text-sm text-zinc-500">右の一覧の「追加」で単語を入れてください</p>
          ) : (
            <ol className="mt-3 flex flex-col gap-2">
              {selectedWords.map((word, i) => (
                <li key={word.id} className="flex items-center gap-2 rounded-xl bg-orange-50 p-2">
                  <span className="w-5 text-center text-sm font-semibold text-orange-500">{i + 1}</span>
                  <Thumbnail src={word.photo_path ?? word.illustration_path} alt={word.english} size="h-10 w-10" />
                  <span className="min-w-0 flex-1 truncate">
                    {word.english} <span className="text-sm text-zinc-500">{word.japanese}</span>
                  </span>
                  <button type="button" onClick={() => move(i, -1)} disabled={i === 0} className={smallButtonClass} aria-label={`${word.english}を上へ`}>
                    ↑
                  </button>
                  <button
                    type="button"
                    onClick={() => move(i, 1)}
                    disabled={i === selectedWords.length - 1}
                    className={smallButtonClass}
                    aria-label={`${word.english}を下へ`}
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
          <input
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="英単語・日本語訳で絞り込み"
            className={`${inputClass} mt-3 text-sm`}
          />
          {unselectedWords.length === 0 ? (
            <p className="mt-3 text-sm text-zinc-400">
              {keyword ? "当てはまる単語がありません" : allWords.length === 0 ? "まだ単語がありません。先に単語管理で登録してください" : "すべての単語が入っています"}
            </p>
          ) : (
            <ul className="mt-3 flex max-h-[28rem] flex-col gap-2 overflow-y-auto">
              {unselectedWords.map((word) => (
                <li key={word.id} className="flex items-center gap-2 rounded-xl border border-zinc-100 p-2">
                  <Thumbnail src={word.photo_path ?? word.illustration_path} alt={word.english} size="h-10 w-10" />
                  <span className="min-w-0 flex-1 truncate">
                    {word.english} <span className="text-sm text-zinc-500">{word.japanese}</span>
                  </span>
                  <button type="button" onClick={() => setSelectedIds((prev) => [...prev, word.id])} className={smallButtonClass}>
                    追加
                  </button>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      <div className="mt-5">
        <ErrorBox error={state.error} details={state.details} />
      </div>
      <div className="mt-5 flex justify-end gap-2">
        <Link href="/admin/courses" className={secondaryButtonClass}>
          キャンセル
        </Link>
        <button type="submit" disabled={pending} className={`${primaryButtonClass} disabled:opacity-60`}>
          {pending ? "保存中…" : course ? "保存する" : "登録する"}
        </button>
      </div>
    </form>
  );
}
