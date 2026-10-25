import { BackLink, PageHeader } from "@/components/admin/ui";
import { UserForm } from "@/components/admin/UserForm";

export default function AdminNewUserPage() {
  return (
    <>
      <BackLink href="/admin/users">ユーザー一覧</BackLink>
      <PageHeader title="ユーザーを作成" />
      <UserForm />
    </>
  );
}
