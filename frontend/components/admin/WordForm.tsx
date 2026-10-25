"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";
import type { Word } from "@/data/mock-data";
import { useAdminStore } from "./AdminStore";
import { BackLink, cardClass, Field, inputClass, NotFoundMessage, PageHeader, primaryButtonClass, secondaryButtonClass } from "./ui";

// /admin/words/[wordId]/edit 用。URL の id から単語を探してフォームに渡す
export function EditWord() {
  const { wordId } = useParams<{ wordId: string }>();
  const { words } = useAdminStore();
  const word = words.find((w) => w.id === Number(wordId));

  if (!word) {
    return <NotFoundMessage backHref="/admin/words">単語が見つかりません</NotFoundMessage>;
  }
  return <WordForm word={word} />;
}

export function WordForm({ word }: { word?: Word }) {
  const router = useRouter();
  const { saveWord } = useAdminStore();

  const [english, setEnglish] = useState(word?.english ?? "");
  const [japanese, setJapanese] = useState(word?.japanese ?? "");
  const [example, setExample] = useState(word?.example ?? "");
  const [photoPath, setPhotoPath] = useState(word?.photo_path ?? "");
  const [illustrationPath, setIllustrationPath] = useState(word?.illustration_path ?? "");
  const [error, setError] = useState("");

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!english.trim() || !japanese.trim()) {
      setError("英単語と日本語訳を入力してください");
      return;
    }
    saveWord({
      id: word?.id ?? null,
      english: english.trim(),
      japanese: japanese.trim(),
      example: example.trim() || null,
      photo_path: photoPath.trim() || null,
      illustration_path: illustrationPath.trim() || null,
    });
    router.push("/admin/words");
  }

  return (
    <form onSubmit={handleSubmit}>
      <BackLink href="/admin/words">単語一覧</BackLink>
      <PageHeader title={word ? "単語を編集" : "単語を登録"} />

      <div className={`${cardClass} flex flex-col gap-4 p-5`}>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="英単語" required>
            <input value={english} onChange={(e) => setEnglish(e.target.value)} maxLength={100} placeholder="dog" className={inputClass} />
          </Field>
          <Field label="日本語訳" required>
            <input value={japanese} onChange={(e) => setJapanese(e.target.value)} maxLength={100} placeholder="犬" className={inputClass} />
          </Field>
        </div>
        <Field label="例文">
          <input value={example} onChange={(e) => setExample(e.target.value)} maxLength={255} placeholder="I have a dog." className={inputClass} />
        </Field>
        <Field label="写真のパス">
          <input
            value={photoPath}
            onChange={(e) => setPhotoPath(e.target.value)}
            maxLength={255}
            placeholder="/images/photo/dog.jpg"
            className={`${inputClass} font-mono text-sm`}
          />
        </Field>
        <Field label="イラストのパス">
          <input
            value={illustrationPath}
            onChange={(e) => setIllustrationPath(e.target.value)}
            maxLength={255}
            placeholder="/images/illust/dog.png"
            className={`${inputClass} font-mono text-sm`}
          />
        </Field>
      </div>

      {error && <p className="mt-4 text-sm text-rose-600">{error}</p>}
      <div className="mt-6 flex justify-end gap-2">
        <Link href="/admin/words" className={secondaryButtonClass}>
          キャンセル
        </Link>
        <button type="submit" className={primaryButtonClass}>
          {word ? "保存する" : "登録する"}
        </button>
      </div>
    </form>
  );
}
