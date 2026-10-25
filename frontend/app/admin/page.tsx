import Link from "next/link";

const MENU = [
  { href: "/admin/courses", title: "コース管理", description: "コースの登録・編集・削除と、出題する単語の並び替え" },
  { href: "/admin/words", title: "単語管理", description: "英単語・日本語訳・例文・画像パスの登録・編集・削除" },
  { href: "/admin/users", title: "ユーザー管理", description: "ユーザーの一覧と、学習の進み具合の確認" },
];

export default function AdminTopPage() {
  return (
    <>
      <h1 className="text-2xl font-bold">管理メニュー</h1>
      <p className="mt-1 text-sm text-zinc-500">
        変更はこの画面の中だけで反映されます（リロードすると元に戻ります）
      </p>
      <ul className="mt-6 flex flex-col gap-3">
        {MENU.map((item) => (
          <li key={item.href}>
            <Link
              href={item.href}
              className="block rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm transition hover:border-orange-300 hover:shadow-md active:scale-[0.99]"
            >
              <p className="text-lg font-semibold">{item.title}</p>
              <p className="mt-1 text-sm text-zinc-500">{item.description}</p>
            </Link>
          </li>
        ))}
      </ul>
    </>
  );
}
