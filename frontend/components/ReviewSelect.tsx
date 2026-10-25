"use client";

import Link from "next/link";
import { useState } from "react";
import { REVIEW_TARGETS, type ReviewTarget } from "@/lib/types";

type Counts = { weak: number; unsure: number };

export type ReviewScope = {
  courseId: number;
  name: string;
  interrupted: boolean;
  counts: Counts;
};

type Props = {
  scopes: ReviewScope[];
  allCounts: Counts | null; // 「すべてのコース」の数。コース固定のときは null
  fixedCourseId: number | null;
};

function countFor(counts: Counts, target: ReviewTarget) {
  return target === "weak" ? counts.weak : target === "unsure" ? counts.unsure : counts.weak + counts.unsure;
}

export function ReviewSelect({ scopes, allCounts, fixedCourseId }: Props) {
  const fixedScope = scopes.find((s) => s.courseId === fixedCourseId);
  const [courseId, setCourseId] = useState<number | "all">(fixedScope?.courseId ?? "all");
  const [target, setTarget] = useState<ReviewTarget>("both");

  const selectedCounts = courseId === "all" ? allCounts : scopes.find((s) => s.courseId === courseId)?.counts;
  const selectedCount = selectedCounts ? countFor(selectedCounts, target) : 0;
  const back = fixedScope ? "course" : "mypage";
  const backHref = fixedScope ? `/courses/${fixedScope.courseId}` : "/mypage";

  return (
    <main className="mx-auto w-full max-w-md flex-1 px-4 py-8">
      <Link href={backHref} className="mb-3 inline-block text-sm text-zinc-500 hover:text-orange-500">
        ← {fixedScope ? "学習に戻る" : "マイページ"}
      </Link>
      <h1 className="text-2xl font-bold">復習する単語を選ぶ</h1>

      <section className="mt-6">
        <h2 className="mb-2 font-semibold">コース</h2>
        {fixedScope ? (
          <p className="rounded-2xl border border-orange-300 bg-white p-4 font-semibold">{fixedScope.name}</p>
        ) : (
          <div className="flex flex-col gap-2">
            {allCounts && (
              <Choice selected={courseId === "all"} onClick={() => setCourseId("all")} label="すべてのコース" count={countFor(allCounts, target)} />
            )}
            {scopes.map((scope) => (
              <Choice
                key={scope.courseId}
                selected={courseId === scope.courseId}
                onClick={() => setCourseId(scope.courseId)}
                label={scope.name}
                badge={scope.interrupted ? "中断中" : undefined}
                count={countFor(scope.counts, target)}
              />
            ))}
          </div>
        )}
      </section>

      <section className="mt-6">
        <h2 className="mb-2 font-semibold">対象</h2>
        <div className="grid grid-cols-3 gap-2">
          {REVIEW_TARGETS.map((t) => (
            <button
              key={t.value}
              type="button"
              onClick={() => setTarget(t.value)}
              className={`rounded-2xl border py-3 transition ${target === t.value ? "border-orange-400 bg-orange-50 font-semibold text-orange-600" : "border-zinc-200 bg-white text-zinc-700 hover:border-orange-200"}`}
            >
              {t.label}
              <span className="block text-xs font-normal text-zinc-500">{selectedCounts ? countFor(selectedCounts, t.value) : 0} 語</span>
            </button>
          ))}
        </div>
      </section>

      {selectedCount > 0 ? (
        <Link
          href={`/review/session?course=${courseId}&target=${target}&back=${back}`}
          className="mt-8 block w-full rounded-2xl bg-orange-500 py-4 text-center font-semibold text-white hover:bg-orange-600"
        >
          {selectedCount} 語を復習する
        </Link>
      ) : (
        <p className="mt-8 rounded-2xl bg-zinc-100 py-4 text-center text-sm text-zinc-500">復習する単語はありません</p>
      )}
    </main>
  );
}

function Choice({
  selected,
  onClick,
  label,
  badge,
  count,
}: {
  selected: boolean;
  onClick: () => void;
  label: string;
  badge?: string;
  count: number;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`flex items-center justify-between rounded-2xl border p-4 text-left transition ${selected ? "border-orange-400 bg-orange-50" : "border-zinc-200 bg-white hover:border-orange-200"}`}
    >
      <span className="flex items-center gap-2">
        <span className={selected ? "font-semibold text-orange-600" : ""}>{label}</span>
        {badge && <span className="rounded-full bg-orange-100 px-2 py-0.5 text-xs font-semibold text-orange-600">{badge}</span>}
      </span>
      <span className="text-sm text-zinc-500">{count} 語</span>
    </button>
  );
}
