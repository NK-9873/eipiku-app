// user-service・learning-service の呼び出し（サーバー側だけで使う）。
// トークンは HttpOnly の Cookie にあり、ここで Authorization: Bearer に付ける。ブラウザは API を直接呼ばない
import "server-only";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { TOKEN_COOKIE } from "./session-cookie";

const SERVICES = {
  user: process.env.USER_API_URL,
  learning: process.env.LEARNING_API_URL,
};

// API がエラーを返したとき。message・details は API の {"error", "details"} をそのまま入れる
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public details: string[] = [],
  ) {
    super(message);
    this.name = "ApiError";
  }
}

type Options = {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown; // JSON にして送る。FormData ならそのまま multipart/form-data で送る
  // ログイン・新規登録のようにトークンが要らない呼び出しは false
  auth?: boolean;
};

export async function apiFetch<T>(service: keyof typeof SERVICES, path: string, options: Options = {}): Promise<T> {
  const baseURL = SERVICES[service];
  if (!baseURL) {
    throw new Error(`${service === "user" ? "USER_API_URL" : "LEARNING_API_URL"} が設定されていません（frontend/.env.local を確認してください）`);
  }

  const isForm = options.body instanceof FormData;
  const headers: Record<string, string> = {};
  if (options.body !== undefined && !isForm) headers["Content-Type"] = "application/json";
  if (options.auth !== false) {
    const token = (await cookies()).get(TOKEN_COOKIE)?.value;
    if (!token) redirect("/login");
    headers.Authorization = `Bearer ${token}`;
  }

  let res: Response;
  try {
    res = await fetch(baseURL + path, {
      method: options.method ?? "GET",
      headers,
      body: options.body === undefined ? undefined : isForm ? (options.body as FormData) : JSON.stringify(options.body),
      cache: "no-store",
    });
  } catch {
    throw new ApiError(503, "サーバーに接続できませんでした。しばらくしてからもう一度お試しください");
  }

  // トークンが無効・期限切れならログインからやり直す（ログイン自体の失敗は呼び出し側でメッセージを出す）
  if (res.status === 401 && options.auth !== false) redirect("/login");

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new ApiError(res.status, body?.error ?? `エラーが発生しました（HTTP ${res.status}）`, body?.details ?? []);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}
