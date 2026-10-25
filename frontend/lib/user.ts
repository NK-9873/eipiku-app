// ログイン中のユーザーの取得（user-service）。サーバー側だけで使う
import "server-only";
import { cache } from "react";
import { apiFetch } from "./api";
import type { User } from "./types";

// ヘッダーと画面の両方で使うので、1回のリクエストの中では1回だけ API を呼ぶ
export const getMe = cache(() => apiFetch<User>("user", "/me"));
