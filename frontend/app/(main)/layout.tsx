import { Suspense } from "react";
import { AppHeader } from "@/components/AppHeader";

export default function MainLayout({ children }: LayoutProps<"/">) {
  return (
    <>
      <Suspense>
        <AppHeader />
      </Suspense>
      {children}
    </>
  );
}
