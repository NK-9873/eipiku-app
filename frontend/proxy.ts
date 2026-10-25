// ログインしていない（トークンの Cookie が無い）ときは、ログイン・新規登録以外の画面をログイン画面へ移す。
// ここでは Cookie があるかだけを見る。トークンが無効・期限切れなら、API が 401 を返した時点で lib/api.ts がログイン画面へ移す
import { NextResponse, type NextRequest } from "next/server";
import { TOKEN_COOKIE } from "@/lib/session-cookie";

const AUTH_PAGES = ["/login", "/signup"];

export function proxy(request: NextRequest) {
  const loggedIn = request.cookies.has(TOKEN_COOKIE);
  const isAuthPage = AUTH_PAGES.includes(request.nextUrl.pathname);

  if (!loggedIn && !isAuthPage) {
    return NextResponse.redirect(new URL("/login", request.url));
  }
  // ログイン済みでもログイン画面は開けるようにしている。
  // Cookie はあってもトークンが無効なとき、ログイン画面とマイページの間で移動を繰り返さないため
  return NextResponse.next();
}

export const config = {
  // 画面だけを対象にする（Next.js の内部ファイルや画像などは除く）
  matcher: ["/((?!_next/static|_next/image|favicon.ico|.*\\.(?:png|jpg|jpeg|webp|svg|ico)$).*)"],
};
