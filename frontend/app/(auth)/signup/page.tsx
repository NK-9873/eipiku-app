import { Suspense } from "react";
import { AuthForm } from "@/components/AuthForm";

export default function SignupPage() {
  // AuthForm は URL のクエリ（?registered=1）を読むので Suspense で包む
  return (
    <Suspense>
      <AuthForm mode="signup" />
    </Suspense>
  );
}
