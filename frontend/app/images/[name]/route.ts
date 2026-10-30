// 単語の画像（GET /images/{ファイル名}）を learning-service から取ってきて返す。
// 本番では learning-service をインターネットに公開しないので、ブラウザは frontend のこの URL から画像を読む。
// next.config の rewrites は next build の時点で転送先が固定されるため、起動時の LEARNING_API_URL を使えるここで転送している
import type { NextRequest } from "next/server";

// ブラウザのキャッシュの確認（304）と、画像の表示に必要なヘッダーだけをやりとりする
const REQUEST_HEADERS = ["if-none-match", "if-modified-since", "range"];
const RESPONSE_HEADERS = [
  "content-type",
  "content-length",
  "cache-control",
  "etag",
  "last-modified",
  "accept-ranges",
  "content-range",
  "x-content-type-options",
];

export async function GET(request: NextRequest, { params }: { params: Promise<{ name: string }> }) {
  const baseURL = process.env.LEARNING_API_URL;
  if (!baseURL) return new Response("LEARNING_API_URL が設定されていません", { status: 500 });

  const { name } = await params;
  const headers = new Headers();
  for (const key of REQUEST_HEADERS) {
    const value = request.headers.get(key);
    if (value) headers.set(key, value);
  }

  let res: Response;
  try {
    res = await fetch(`${baseURL}/images/${encodeURIComponent(name)}`, { headers, cache: "no-store" });
  } catch {
    return new Response("画像を取得できませんでした", { status: 502 });
  }

  const resHeaders = new Headers();
  for (const key of RESPONSE_HEADERS) {
    const value = res.headers.get(key);
    if (value) resHeaders.set(key, value);
  }
  return new Response(res.body, { status: res.status, headers: resHeaders });
}
