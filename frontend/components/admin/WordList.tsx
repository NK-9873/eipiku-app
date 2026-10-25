import Link from "next/link";
import { deleteWord } from "@/lib/admin-actions";
import type { Word } from "@/lib/types";
import { DeleteButton } from "./DeleteButton";
import { Thumbnail } from "./Thumbnail";
import { cardClass, PageHeader, primaryButtonClass, SavedNotice, smallButtonClass } from "./ui";

// 単語一覧（GET /words）
export function WordList({ words, saved }: { words: Word[]; saved?: string }) {
  return (
    <>
      <PageHeader
        title="単語管理"
        action={
          <Link href="/admin/words/new" className={primaryButtonClass}>
            ＋ 単語を登録
          </Link>
        }
      />
      {saved && <SavedNotice>「{saved}」を保存しました</SavedNotice>}

      {words.length === 0 ? (
        <p className={`${cardClass} p-8 text-center text-zinc-500`}>まだ単語がありません。「＋ 単語を登録」から追加してください</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {words.map((word) => (
            <li key={word.id} className={`${cardClass} flex items-start gap-4 p-4`}>
              <div className="flex gap-1.5">
                <Thumbnail src={word.photo_path} alt={`${word.english}の写真`} size="h-16 w-16" />
                <Thumbnail src={word.illustration_path} alt={`${word.english}のイラスト`} size="h-16 w-16" />
              </div>
              <div className="min-w-0 flex-1">
                <p>
                  <span className="text-lg font-semibold">{word.english}</span>
                  <span className="ml-2 text-orange-500">{word.japanese}</span>
                </p>
                {word.example && <p className="mt-1 text-sm text-zinc-600">{word.example}</p>}
                {word.example_japanese && <p className="text-sm text-zinc-400">{word.example_japanese}</p>}
              </div>
              <div className="flex shrink-0 items-start gap-1">
                <Link href={`/admin/words/${word.id}/edit`} className={smallButtonClass}>
                  編集
                </Link>
                <DeleteButton
                  confirmMessage={`単語「${word.english}」を削除しますか？\nこの単語を含むコースからも外れ、学習の記録も消えます。元に戻せません。`}
                  action={deleteWord.bind(null, word.id)}
                />
              </div>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}
