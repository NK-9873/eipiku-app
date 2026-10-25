"use client";

import { useImageError } from "../useImageError";

// 単語の画像の小さな表示。画像が無い・読み込めない（古いパスなど）ときは灰色の枠
export function Thumbnail({ src, alt, size = "h-12 w-12" }: { src: string | null; alt: string; size?: string }) {
  const [failed, imageRef, onImageError] = useImageError(src);

  if (!src || failed) {
    return (
      <span className={`${size} flex shrink-0 items-center justify-center rounded-lg bg-zinc-100 text-[10px] text-zinc-400`}>
        画像なし
      </span>
    );
  }
  return (
    // eslint-disable-next-line @next/next/no-img-element -- learning-service の画像を、最適化せずそのまま表示する
    <img
      ref={imageRef}
      src={src}
      alt={alt}
      onError={onImageError}
      className={`${size} shrink-0 rounded-lg border border-zinc-200 bg-zinc-50 object-cover`}
    />
  );
}
