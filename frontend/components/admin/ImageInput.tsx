"use client";

import { useEffect, useRef, useState, type ClipboardEvent, type DragEvent } from "react";
import { IMAGE_ACCEPT, MAX_IMAGE_BYTES } from "@/lib/types";
import { useImageError } from "../useImageError";
import { secondaryButtonClass, smallButtonClass } from "./ui";

// 画像の欄の状態。url は登録済みの画像、file はこれから登録する画像（保存するときにアップロードする）
export type ImageValue = { kind: "none" } | { kind: "url"; url: string } | { kind: "file"; file: File; previewUrl: string };

// 選んだファイルを確かめる。問題があれば、そのまま画面に出せるメッセージを返す
function checkFile(file: File): string | null {
  if (!IMAGE_ACCEPT.includes(file.type)) {
    return `「${file.name || "貼り付けたデータ"}」は使えない形式です。JPEG・PNG・WebP の画像を選んでください`;
  }
  if (file.size > MAX_IMAGE_BYTES) {
    return `画像が大きすぎます（${(file.size / 1024 / 1024).toFixed(1)}MB）。5MB以下の画像を選んでください`;
  }
  return null;
}

type Props = {
  label: string;
  value: ImageValue;
  onChange: (value: ImageValue) => void;
  error?: string; // 保存したときに API から返ったエラー
};

// 画像の登録欄。貼り付け（Ctrl/⌘+V）・ドラッグ&ドロップ・ファイル選択で画像を選び、プレビューを出す
export function ImageInput({ label, value, onChange, error }: Props) {
  const fileInput = useRef<HTMLInputElement>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [dragging, setDragging] = useState(false);

  // プレビュー用に作った URL は、使わなくなったら解放する
  const previewUrl = value.kind === "file" ? value.previewUrl : null;
  useEffect(() => () => void (previewUrl && URL.revokeObjectURL(previewUrl)), [previewUrl]);

  function choose(file: File | undefined | null) {
    if (!file) {
      setMessage("画像が見つかりませんでした。画像をコピーしてから貼り付けてください");
      return;
    }
    const problem = checkFile(file);
    if (problem) {
      setMessage(problem);
      return;
    }
    setMessage(null);
    onChange({ kind: "file", file, previewUrl: URL.createObjectURL(file) });
  }

  function handlePaste(e: ClipboardEvent) {
    const file = [...e.clipboardData.items].find((item) => item.kind === "file")?.getAsFile();
    e.preventDefault();
    choose(file);
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    setDragging(false);
    choose(e.dataTransfer.files[0]);
  }

  const src = value.kind === "file" ? value.previewUrl : value.kind === "url" ? value.url : null;
  const [broken, imageRef, onImageError] = useImageError(src); // 登録済みの URL が表示できないとき
  const shownError = message ?? error;

  return (
    <div>
      <p className="mb-1 text-sm font-medium text-zinc-700">{label}</p>
      <div
        tabIndex={0}
        role="button"
        aria-label={`${label}の画像。クリックしてから貼り付け、またはドラッグ&ドロップ`}
        onPaste={handlePaste}
        onDragOver={(e) => {
          e.preventDefault();
          setDragging(true);
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={handleDrop}
        className={`group flex aspect-square w-full flex-col items-center justify-center overflow-hidden rounded-2xl border-2 border-dashed text-center outline-none transition focus:border-orange-400 focus:bg-orange-50 ${
          dragging ? "border-orange-400 bg-orange-50" : shownError ? "border-rose-300 bg-rose-50" : "border-zinc-300 bg-zinc-50"
        }`}
      >
        {src && !broken ? (
          // eslint-disable-next-line @next/next/no-img-element -- 選んだ画像（ブラウザ内の URL）や learning-service の画像をそのまま表示する
          <img ref={imageRef} src={src} alt={`${label}のプレビュー`} className="h-full w-full object-contain" onError={onImageError} />
        ) : (
          <div className="px-4 text-sm text-zinc-500">
            {broken ? (
              <p className="text-rose-600">登録されている画像を表示できません。新しい画像を登録してください</p>
            ) : (
              <>
                <p className="font-semibold text-zinc-600">ここに画像をドラッグ</p>
                <p className="mt-1">または、ここをクリックしてから</p>
                <p>Ctrl+V（Mac は ⌘+V）で貼り付け</p>
              </>
            )}
          </div>
        )}
      </div>
      <p className="mt-1 text-xs text-zinc-400">JPEG・PNG・WebP、5MBまで{value.kind === "file" && "（保存すると登録されます）"}</p>

      {shownError && (
        <p role="alert" className="mt-1 text-sm text-rose-600">
          {shownError}
        </p>
      )}

      <div className="mt-2 flex gap-2">
        <button type="button" onClick={() => fileInput.current?.click()} className={`${secondaryButtonClass} px-3 py-1.5 text-sm`}>
          {src ? "画像を差し替える" : "ファイルを選ぶ"}
        </button>
        {value.kind !== "none" && (
          <button
            type="button"
            onClick={() => {
              setMessage(null);
              onChange({ kind: "none" });
            }}
            className={`${smallButtonClass} text-rose-600`}
          >
            画像を外す
          </button>
        )}
      </div>
      <input
        ref={fileInput}
        type="file"
        accept={IMAGE_ACCEPT.join(",")}
        className="hidden"
        onChange={(e) => {
          choose(e.target.files?.[0]);
          e.target.value = ""; // 同じファイルをもう一度選んでも反応するように
        }}
      />
    </div>
  );
}
