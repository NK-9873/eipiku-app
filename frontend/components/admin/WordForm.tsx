"use client";

import Link from "next/link";
import { useState, useTransition, type FormEvent } from "react";
import { saveWord } from "@/lib/admin-actions";
import type { FormState, Word } from "@/lib/types";
import { ImageInput, type ImageValue } from "./ImageInput";
import { BackLink, cardClass, ErrorBox, Field, inputClass, PageHeader, primaryButtonClass, secondaryButtonClass } from "./ui";

function initialImage(url: string | null | undefined): ImageValue {
  return url ? { kind: "url", url } : { kind: "none" };
}

// 画像の欄の値を、Server Action に送る形にする（新しいファイルか、今の URL か）
function appendImage(formData: FormData, key: "photo" | "illustration", value: ImageValue) {
  if (value.kind === "file") formData.append(`${key}_file`, value.file);
  if (value.kind === "url") formData.append(`${key}_url`, value.url);
}

// 単語の登録・編集。写真とイラストは保存するときに一緒にアップロードする
export function WordForm({ word }: { word?: Word }) {
  const [photo, setPhoto] = useState<ImageValue>(initialImage(word?.photo_path));
  const [illustration, setIllustration] = useState<ImageValue>(initialImage(word?.illustration_path));
  const [state, setState] = useState<FormState>({});
  const [pending, startTransition] = useTransition();

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    appendImage(formData, "photo", photo);
    appendImage(formData, "illustration", illustration);
    startTransition(async () => {
      // 成功したら一覧に移る（saveWord の中で redirect）。戻ってきたらエラー
      setState(await saveWord(word?.id ?? null, formData));
    });
  }

  // 画像のエラー（「写真：…」のように始まる）は、その画像の欄の下にも出す
  const imageError = (label: string) => (state.error?.startsWith(`${label}：`) ? state.error.slice(label.length + 1) : undefined);

  return (
    <form onSubmit={handleSubmit}>
      <BackLink href="/admin/words">単語一覧</BackLink>
      <PageHeader title={word ? `単語を編集：${word.english}` : "単語を登録"} />

      <div className={`${cardClass} flex flex-col gap-4 p-5`}>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="英単語" required>
            <input name="english" defaultValue={word?.english} required maxLength={100} placeholder="例：dog" className={inputClass} />
          </Field>
          <Field label="日本語訳" required>
            <input name="japanese" defaultValue={word?.japanese} required maxLength={100} placeholder="例：犬" className={inputClass} />
          </Field>
        </div>
        <Field label="例文">
          <input name="example" defaultValue={word?.example ?? ""} maxLength={255} placeholder="例：I have a dog." className={inputClass} />
        </Field>
        <Field label="例文の日本語訳">
          <input
            name="example_japanese"
            defaultValue={word?.example_japanese ?? ""}
            maxLength={255}
            placeholder="例：私は犬を飼っています。"
            className={inputClass}
          />
        </Field>
      </div>

      <div className={`${cardClass} mt-5 p-5`}>
        <h2 className="font-semibold">画像</h2>
        <p className="mt-1 text-sm text-zinc-500">
          学習する人が設定で選んだ方（写真かイラスト）が表示されます。どちらか一方だけでも登録できます。
        </p>
        <div className="mt-4 grid gap-5 sm:grid-cols-2">
          <ImageInput label="写真" value={photo} onChange={setPhoto} error={imageError("写真")} />
          <ImageInput label="イラスト" value={illustration} onChange={setIllustration} error={imageError("イラスト")} />
        </div>
      </div>

      <div className="mt-5">
        <ErrorBox error={state.error} details={state.details} />
      </div>
      <div className="mt-5 flex justify-end gap-2">
        <Link href="/admin/words" className={secondaryButtonClass}>
          キャンセル
        </Link>
        <button type="submit" disabled={pending} className={`${primaryButtonClass} disabled:opacity-60`}>
          {pending ? "保存中…" : word ? "保存する" : "登録する"}
        </button>
      </div>
    </form>
  );
}
