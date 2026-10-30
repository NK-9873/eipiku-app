// ログイン中のトークン（JWT）を入れる Cookie の名前。proxy.ts からも読むので、ほかの処理とは分けている
export const TOKEN_COOKIE = "eipiku_token";

// Cookie に Secure 属性（https のときだけ送る）を付けるか。COOKIE_SECURE=true / false で指定する。
// 未設定なら本番（next build したもの）だけ付ける。https なしの本番では false にしないとログインできない
export function isSecureCookie(): boolean {
  const value = process.env.COOKIE_SECURE;
  if (value === "true") return true;
  if (value === "false") return false;
  return process.env.NODE_ENV === "production";
}
