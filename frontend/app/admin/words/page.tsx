import { Suspense } from "react";
import { WordList } from "@/components/admin/WordList";
import { Loading } from "@/components/admin/ui";
import { getWords } from "@/lib/admin";

export default function AdminWordsPage({ searchParams }: PageProps<"/admin/words">) {
  return (
    <Suspense fallback={<Loading />}>
      <Words searchParams={searchParams} />
    </Suspense>
  );
}

async function Words({ searchParams }: Pick<PageProps<"/admin/words">, "searchParams">) {
  const [{ saved }, words] = await Promise.all([searchParams, getWords()]);
  return <WordList words={words} saved={typeof saved === "string" ? saved : undefined} />;
}
