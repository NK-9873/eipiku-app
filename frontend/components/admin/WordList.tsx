"use client";

import Link from "next/link";
import { useAdminStore } from "./AdminStore";
import { cardClass, dangerButtonClass, PageHeader, primaryButtonClass, smallButtonClass } from "./ui";

export function WordList() {
  const { words, courses, courseWords, deleteWord } = useAdminStore();

  function coursesOf(wordId: number) {
    return courseWords
      .filter((cw) => cw.word_id === wordId)
      .map((cw) => courses.find((c) => c.id === cw.course_id)?.name)
      .filter((name) => name !== undefined);
  }

  function handleDelete(wordId: number, english: string) {
    const used = coursesOf(wordId);
    const note = used.length > 0 ? `\nコース（${used.join("、")}）からも外されます。` : "";
    if (window.confirm(`単語「${english}」を削除しますか？${note}`)) deleteWord(wordId);
  }

  return (
    <>
      <PageHeader
        title="単語管理"
        action={
          <Link href="/admin/words/new" className={primaryButtonClass}>
            ＋ 新規登録
          </Link>
        }
      />

      {words.length === 0 ? (
        <p className={`${cardClass} p-8 text-center text-zinc-500`}>単語がありません</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {words.map((word) => {
            const used = coursesOf(word.id);
            return (
              <li key={word.id} className={`${cardClass} p-5`}>
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p>
                      <span className="text-lg font-semibold">{word.english}</span>
                      <span className="ml-2 text-orange-500">{word.japanese}</span>
                    </p>
                    {word.example && <p className="mt-1 text-sm text-zinc-600">{word.example}</p>}
                  </div>
                  <div className="flex shrink-0 gap-1">
                    <Link href={`/admin/words/${word.id}/edit`} className={smallButtonClass}>
                      編集
                    </Link>
                    <button type="button" onClick={() => handleDelete(word.id, word.english)} className={dangerButtonClass}>
                      削除
                    </button>
                  </div>
                </div>
                <dl className="mt-3 grid grid-cols-[auto_1fr] gap-x-3 gap-y-0.5 text-xs text-zinc-500">
                  <dt>写真</dt>
                  <dd className="truncate font-mono">{word.photo_path ?? "—"}</dd>
                  <dt>イラスト</dt>
                  <dd className="truncate font-mono">{word.illustration_path ?? "—"}</dd>
                  <dt>コース</dt>
                  <dd>{used.length > 0 ? used.join("、") : "未使用"}</dd>
                </dl>
              </li>
            );
          })}
        </ul>
      )}
    </>
  );
}
