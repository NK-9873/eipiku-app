// 管理画面のデータ取得（サーバー側だけで使う）。どの API も管理者でなければ 403 になる
import "server-only";
import { redirect } from "next/navigation";
import { apiFetch, ApiError } from "./api";
import type { CourseWord, User, UserProgress, Word } from "./types";
import { getMe } from "./user";

// 管理者でなければトップへ移す（未ログイン・トークンが無効なら getMe の中でログイン画面へ移る）
export async function requireAdmin(): Promise<User> {
  const me = await getMe();
  if (me.role !== "admin") redirect("/");
  return me;
}

export function getWords() {
  return apiFetch<Word[]>("learning", "/words");
}

export function getWord(id: number) {
  return apiFetch<Word>("learning", `/words/${id}`);
}

export function getUsers() {
  return apiFetch<User[]>("user", "/users");
}

export function getUser(id: number) {
  return apiFetch<User>("user", `/users/${id}`);
}

export function getUserProgress(id: number) {
  return apiFetch<UserProgress>("learning", `/admin/users/${id}/progress`);
}

// 見つからない（404）・ID が不正（400）なら null にする。それ以外のエラーはそのまま投げる
export async function orNotFound<T>(promise: Promise<T>): Promise<T | null> {
  try {
    return await promise;
  } catch (e) {
    if (e instanceof ApiError && (e.status === 404 || e.status === 400)) return null;
    throw e;
  }
}

export type { CourseWord };
