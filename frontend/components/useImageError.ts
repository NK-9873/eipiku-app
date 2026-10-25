"use client";

import { useCallback, useState } from "react";

// 画像を読み込めなかったか（古いパスや削除された画像など）を返す。
// 最初の表示（React の準備ができる前）に読み込みに失敗した画像は onError が呼ばれないので、
// img が画面に付いたとき（ref）にも読み込み済みかどうかを確かめる
export function useImageError(src: string | null | undefined) {
  const [failedSrc, setFailedSrc] = useState<string | null>(null);

  const ref = useCallback(
    (img: HTMLImageElement | null) => {
      if (src && img && img.complete && img.naturalWidth === 0) setFailedSrc(src);
    },
    [src],
  );

  // [読み込めなかったか, img に付ける ref, img の onError]
  return [!!src && failedSrc === src, ref, () => setFailedSrc(src ?? null)] as const;
}
