import { Suspense } from "react";
import { WordForm } from "@/components/admin/WordForm";
import { Loading, NotFoundMessage } from "@/components/admin/ui";
import { getWord, orNotFound } from "@/lib/admin";

export default function AdminEditWordPage({ params }: PageProps<"/admin/words/[wordId]/edit">) {
  return (
    <Suspense fallback={<Loading />}>
      <EditWord params={params} />
    </Suspense>
  );
}

async function EditWord({ params }: Pick<PageProps<"/admin/words/[wordId]/edit">, "params">) {
  const word = await orNotFound(getWord(Number((await params).wordId)));
  if (!word) return <NotFoundMessage backHref="/admin/words">単語が見つかりません。削除された可能性があります</NotFoundMessage>;
  return <WordForm key={word.id} word={word} />;
}
