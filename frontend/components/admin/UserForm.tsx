"use client";

import Link from "next/link";
import { useActionState } from "react";
import { createUser, updateUser } from "@/lib/admin-actions";
import type { FormState, User } from "@/lib/types";
import { cardClass, ErrorBox, Field, inputClass, primaryButtonClass, secondaryButtonClass } from "./ui";

type Props = {
  user?: User; // 無ければ作成
  isSelf?: boolean; // ログイン中の管理者自身を編集しているとき
};

// ユーザーの作成（POST /users）・変更（PUT /users/{id}）
export function UserForm({ user, isSelf = false }: Props) {
  const action = user ? updateUser.bind(null, user.id) : createUser;
  const [state, formAction, pending] = useActionState<FormState, FormData>(action, {});

  return (
    <form action={formAction} className={`${cardClass} flex flex-col gap-4 p-5`}>
      <h2 className="font-semibold">{user ? "アカウントを変更" : "アカウントの情報"}</h2>
      <Field label="名前" required>
        <input name="name" defaultValue={state.values?.name ?? user?.name} required maxLength={50} className={inputClass} />
      </Field>
      <Field label="メールアドレス" required>
        <input
          name="email"
          type="email"
          defaultValue={state.values?.email ?? user?.email}
          required
          maxLength={255}
          autoComplete="off"
          className={inputClass}
        />
      </Field>
      <Field label={user ? "新しいパスワード" : "パスワード"} required={!user}>
        <input
          name="password"
          type="password"
          required={!user}
          minLength={8}
          autoComplete="new-password"
          placeholder={user ? "変更するときだけ入力" : ""}
          className={inputClass}
        />
        <span className="mt-1 block text-xs text-zinc-400">
          8文字以上{user && "。空のままなら、今のパスワードのまま変わりません"}
        </span>
      </Field>

      <fieldset>
        <legend className="mb-1 text-sm font-medium text-zinc-700">権限</legend>
        <div className="grid grid-cols-2 gap-2">
          {[
            { value: "user", label: "一般", description: "学習だけできる" },
            { value: "admin", label: "管理者", description: "この管理画面も使える" },
          ].map((r) => (
            <label
              key={r.value}
              className="flex cursor-pointer flex-col rounded-xl border border-zinc-200 p-3 has-[:checked]:border-orange-400 has-[:checked]:bg-orange-50"
            >
              <span className="flex items-center gap-2 font-semibold">
                <input type="radio" name="role" value={r.value} defaultChecked={(user?.role ?? "user") === r.value} />
                {r.label}
              </span>
              <span className="mt-0.5 pl-5 text-xs text-zinc-500">{r.description}</span>
            </label>
          ))}
        </div>
        {user && (
          <p className="mt-2 rounded-lg bg-amber-50 px-3 py-2 text-xs text-amber-800">
            ⚠ 権限の変更は、本人が<strong>次にログインしたとき</strong>に反映されます（ログイン中の人は、ログインし直すまで今の権限のままです。最長24時間）。
            {isSelf && " ご自身の権限を「一般」にすると、次にログインしたときからこの管理画面を開けなくなります。"}
          </p>
        )}
      </fieldset>

      <ErrorBox error={state.error} details={state.details} />

      <div className="flex justify-end gap-2">
        <Link href="/admin/users" className={secondaryButtonClass}>
          {user ? "一覧に戻る" : "キャンセル"}
        </Link>
        <button type="submit" disabled={pending} className={`${primaryButtonClass} disabled:opacity-60`}>
          {pending ? "保存中…" : user ? "変更を保存する" : "作成する"}
        </button>
      </div>
    </form>
  );
}
