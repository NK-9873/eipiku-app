"use client";

import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";
import type { FormState } from "@/lib/types";
import { dangerButtonClass } from "./ui";

// 確認してから削除し、一覧を読み直す。失敗したら理由をその場に出す
export function DeleteButton({ confirmMessage, action }: { confirmMessage: string; action: () => Promise<FormState> }) {
  const router = useRouter();
  const [pending, startTransition] = useTransition();
  const [error, setError] = useState<string | undefined>();

  function handleClick() {
    if (!window.confirm(confirmMessage)) return;
    setError(undefined);
    startTransition(async () => {
      const result = await action();
      if (result.error) {
        setError(result.error);
        return;
      }
      router.refresh();
    });
  }

  return (
    <span className="flex flex-col items-end">
      <button type="button" onClick={handleClick} disabled={pending} className={`${dangerButtonClass} disabled:opacity-50`}>
        {pending ? "削除中…" : "削除"}
      </button>
      {error && (
        <span role="alert" className="mt-1 max-w-56 text-right text-xs text-rose-600">
          {error}
        </span>
      )}
    </span>
  );
}
