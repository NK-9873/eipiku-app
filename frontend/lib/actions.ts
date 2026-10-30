"use server";

// 画面（ブラウザ）から呼ぶ処理。Server Action としてサーバーで動き、トークンの Cookie を読んで API を呼ぶ。
// 直接 POST されても動くので、どの処理も API 側でトークンを確かめている（無効なら apiFetch がログイン画面へ移す）
import { cookies } from "next/headers";
import { redirect, unstable_rethrow } from "next/navigation";
import { apiFetch, ApiError } from "./api";
import { isSecureCookie, TOKEN_COOKIE } from "./session-cookie";
import type { FormState, ImageType, Mastery, User } from "./types";

function toFormState(e: unknown): FormState {
  unstable_rethrow(e); // redirect などは画面のエラーにせず、そのまま Next.js に任せる
  if (e instanceof ApiError) return { error: e.message, details: e.details };
  console.error(e);
  return { error: "エラーが発生しました。しばらくしてからもう一度お試しください" };
}

// ログイン（user-service の POST /login）。成功したらトークンを HttpOnly の Cookie に入れてマイページへ
export async function login(_prev: FormState, formData: FormData): Promise<FormState> {
  const email = String(formData.get("email") ?? "");
  let result: { token: string; expires_at: string };
  try {
    result = await apiFetch("user", "/login", {
      method: "POST",
      auth: false,
      body: { email, password: String(formData.get("password") ?? "") },
    });
  } catch (e) {
    return { ...toFormState(e), values: { email } };
  }
  (await cookies()).set(TOKEN_COOKIE, result.token, {
    httpOnly: true, // JavaScript から読めないようにする
    sameSite: "lax",
    secure: isSecureCookie(), // COOKIE_SECURE で切り替える（https なしで動かすときは false）
    path: "/",
    expires: new Date(result.expires_at), // トークンと同じ24時間で消える
  });
  redirect("/mypage");
}

// 新規登録（user-service の POST /signup）。成功したらログイン画面へ
export async function signup(_prev: FormState, formData: FormData): Promise<FormState> {
  const name = String(formData.get("name") ?? "");
  const email = String(formData.get("email") ?? "");
  try {
    await apiFetch("user", "/signup", {
      method: "POST",
      auth: false,
      body: { name, email, password: String(formData.get("password") ?? "") },
    });
  } catch (e) {
    return { ...toFormState(e), values: { name, email } };
  }
  redirect("/login?registered=1");
}

export async function logout() {
  (await cookies()).delete(TOKEN_COOKIE);
  redirect("/login");
}

// 単語学習で1語評価したとき。習得度と中断位置（次に出す単語の位置）をまとめて保存する
export async function recordStudy(input: { courseId: number; wordId: number; mastery: Mastery; nextPosition: number }) {
  await Promise.all([
    saveWordProgress(input.wordId, input.mastery),
    saveCourseProgress(input.courseId, input.nextPosition),
  ]);
}

// 習得度の保存（PUT /me/word-progresses/{wordId}）
export async function saveWordProgress(wordId: number, mastery: Mastery) {
  await apiFetch("learning", `/me/word-progresses/${wordId}`, { method: "PUT", body: { mastery } });
}

// 中断位置の保存（PUT /me/course-progresses/{courseId}）
export async function saveCourseProgress(courseId: number, currentPosition: number) {
  await apiFetch("learning", `/me/course-progresses/${courseId}`, {
    method: "PUT",
    body: { current_position: currentPosition },
  });
}

// 画像タイプの変更（user-service の PUT /me/settings）
export async function updateImageType(imageType: ImageType): Promise<User> {
  return apiFetch<User>("user", "/me/settings", { method: "PUT", body: { image_type: imageType } });
}
