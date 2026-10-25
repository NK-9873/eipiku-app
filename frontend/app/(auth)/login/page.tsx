import { Suspense } from "react";
import { AuthForm } from "@/components/AuthForm";

export default function LoginPage() {
  // AuthForm は URL のクエリ（?registered=1）を読むので Suspense で包む
  return (
    <Suspense>
      <AuthForm mode="login" />
    </Suspense>
  );
}
