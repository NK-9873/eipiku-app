"use server";

// 管理画面から呼ぶ処理（Server Action）。どれも最初に管理者か確かめる（API 側でも確かめている）。
// エラーは、プログラマでない人にも分かる言葉に直して画面に返す
import { revalidatePath } from "next/cache";
import { redirect, unstable_rethrow } from "next/navigation";
import { requireAdmin } from "./admin";
import { apiFetch, ApiError } from "./api";
import { IMAGE_ACCEPT, MAX_IMAGE_BYTES, type Course, type FormState, type Role, type User, type Word } from "./types";

// API のメッセージに出てくる項目名（英語）を、画面の項目名に置き換える。長い名前から順に置き換える
const FIELD_LABELS: [string, string][] = [
  ["example_japanese", "例文の日本語訳"],
  ["illustration_path", "イラスト"],
  ["photo_path", "写真"],
  ["word_ids", "単語"],
  ["description", "説明"],
  ["japanese", "日本語訳"],
  ["english", "英単語"],
  ["example", "例文"],
  ["password", "パスワード"],
  ["email", "メールアドレス"],
  ["name", "名前"],
  ["role", "権限"],
];

// 文の途中に出てくる言い回し
const PHRASES: [string, string][] = [
  ["role を user に", "権限を「一般」に"],
  ["admin / user", "管理者 / 一般"],
];

function friendly(text: string) {
  for (const [field, label] of FIELD_LABELS) {
    if (text.startsWith(field + " ")) {
      text = label + text.slice(field.length + 1); // 「name は必須です」→「名前は必須です」
      break;
    }
  }
  for (const [from, to] of PHRASES) text = text.replace(from, to);
  return text;
}

function toFormState(e: unknown, context?: string): FormState {
  unstable_rethrow(e); // redirect などはそのまま Next.js に任せる
  if (e instanceof ApiError) {
    const prefix = context ? `${context}：` : "";
    switch (e.status) {
      case 403:
        return { error: `${prefix}この操作を行う権限がありません。管理者のアカウントでログインし直してください` };
      case 404:
        return { error: `${prefix}対象が見つかりません。ほかの人が削除した可能性があります。一覧に戻ってやり直してください` };
      case 503:
        return { error: `${prefix}${e.message}` };
      default:
        return { error: prefix + friendly(e.message), details: e.details.map(friendly) };
    }
  }
  console.error(e);
  return { error: "エラーが発生しました。しばらくしてからもう一度お試しください" };
}

// フォームの文字の項目。空なら null（任意項目）
function text(formData: FormData, key: string) {
  const v = formData.get(key);
  return typeof v === "string" && v.trim() !== "" ? v.trim() : null;
}

// ---- 画像 ----

// 画像の欄の値を決める。新しいファイルがあればアップロードして URL を、
// なければ今の URL（削除したなら null）を返す
async function resolveImage(formData: FormData, key: "photo" | "illustration", label: string): Promise<string | null> {
  const file = formData.get(`${key}_file`);
  if (file instanceof File && file.size > 0) {
    if (!IMAGE_ACCEPT.includes(file.type)) {
      throw new ApiError(400, `${label}：JPEG・PNG・WebP の画像を選んでください`);
    }
    if (file.size > MAX_IMAGE_BYTES) {
      throw new ApiError(400, `${label}：画像が大きすぎます。5MB以下にしてください`);
    }
    const body = new FormData();
    body.append("file", file, file.name || `${key}.bin`);
    try {
      const uploaded = await apiFetch<{ url: string }>("learning", "/images", { method: "POST", body });
      return uploaded.url;
    } catch (e) {
      if (e instanceof ApiError) throw new ApiError(e.status, `${label}：${e.details[0] ?? e.message}`);
      throw e;
    }
  }
  return text(formData, `${key}_url`);
}

// ---- 単語 ----

// 単語の登録・更新。写真とイラストはこの中で POST /images にアップロードし、返った URL を保存する
export async function saveWord(wordId: number | null, formData: FormData): Promise<FormState> {
  await requireAdmin();
  try {
    const [photo, illustration] = [
      await resolveImage(formData, "photo", "写真"),
      await resolveImage(formData, "illustration", "イラスト"),
    ];
    const body = {
      english: String(formData.get("english") ?? ""),
      japanese: String(formData.get("japanese") ?? ""),
      example: text(formData, "example"),
      example_japanese: text(formData, "example_japanese"),
      photo_path: photo,
      illustration_path: illustration,
    };
    if (wordId === null) {
      await apiFetch<Word>("learning", "/words", { method: "POST", body });
    } else {
      await apiFetch<Word>("learning", `/words/${wordId}`, { method: "PUT", body });
    }
  } catch (e) {
    return toFormState(e);
  }
  revalidatePath("/admin/words");
  redirect(`/admin/words?saved=${encodeURIComponent(String(formData.get("english") ?? ""))}`);
}

export async function deleteWord(wordId: number): Promise<FormState> {
  await requireAdmin();
  try {
    await apiFetch("learning", `/words/${wordId}`, { method: "DELETE" });
  } catch (e) {
    return toFormState(e, "削除できませんでした");
  }
  revalidatePath("/admin/words");
  return {};
}

// ---- コース ----

// コースの登録・更新と、含める単語（出題順）の保存（PUT /courses/{id}/words）
export async function saveCourse(
  courseId: number | null,
  input: { name: string; description: string; wordIds: number[] },
): Promise<FormState & { courseId?: number }> {
  await requireAdmin();
  const body = { name: input.name, description: input.description.trim() || null };
  let id = courseId;
  try {
    if (id === null) {
      id = (await apiFetch<Course>("learning", "/courses", { method: "POST", body })).id;
    } else {
      await apiFetch<Course>("learning", `/courses/${id}`, { method: "PUT", body });
    }
  } catch (e) {
    return toFormState(e);
  }
  try {
    await apiFetch("learning", `/courses/${id}/words`, { method: "PUT", body: { word_ids: input.wordIds } });
  } catch (e) {
    // コース自体は保存できているので、編集画面でやり直せるよう id を返す
    return { ...toFormState(e, "コースは保存しましたが、単語の保存に失敗しました"), courseId: id };
  }
  revalidatePath("/admin/courses");
  redirect(`/admin/courses?saved=${encodeURIComponent(input.name)}`);
}

export async function deleteCourse(courseId: number): Promise<FormState> {
  await requireAdmin();
  try {
    await apiFetch("learning", `/courses/${courseId}`, { method: "DELETE" });
  } catch (e) {
    return toFormState(e, "削除できませんでした");
  }
  revalidatePath("/admin/courses");
  return {};
}

// ---- ユーザー ----

function userBody(formData: FormData) {
  return {
    name: String(formData.get("name") ?? ""),
    email: String(formData.get("email") ?? ""),
    password: String(formData.get("password") ?? ""),
    role: (String(formData.get("role") ?? "user") as Role) || "user",
  };
}

export async function createUser(_prev: FormState, formData: FormData): Promise<FormState> {
  await requireAdmin();
  const body = userBody(formData);
  let user: User;
  try {
    user = await apiFetch<User>("user", "/users", { method: "POST", body });
  } catch (e) {
    return { ...toFormState(e), values: { name: body.name, email: body.email } };
  }
  revalidatePath("/admin/users");
  redirect(`/admin/users/${user.id}?saved=created`);
}

// パスワードは入力したときだけ変更する（空なら API でも変更しない）
export async function updateUser(userId: number, _prev: FormState, formData: FormData): Promise<FormState> {
  await requireAdmin();
  const body = userBody(formData);
  try {
    await apiFetch<User>("user", `/users/${userId}`, { method: "PUT", body });
  } catch (e) {
    return { ...toFormState(e), values: { name: body.name, email: body.email } };
  }
  revalidatePath("/admin/users");
  redirect(`/admin/users/${userId}?saved=updated`);
}
