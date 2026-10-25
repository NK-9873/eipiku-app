// ログイン中のユーザーの取得と設定の保存。
// ログイン機能ができるまでは、モックの users の id=1 に固定している。
import { users, type User } from "@/data/mock-data";
import { readStore, writeStore } from "./local-store";

export type { User };
export type ImageType = User["image_type"];

export const CURRENT_USER_ID = 1;

export const IMAGE_TYPE_LABELS: Record<ImageType, string> = {
  photo: "写真",
  illustration: "イラスト",
};

// users.image_type だけ localStorage で上書きできるようにしている
type StoredUserSettings = Record<number, { image_type: ImageType }>;
const SETTINGS_KEY = "user_settings";

export async function getCurrentUser(): Promise<User> {
  const user = users.find((u) => u.id === CURRENT_USER_ID);
  if (!user) throw new Error(`id=${CURRENT_USER_ID} のユーザーがモックデータにありません`);
  const settings = readStore<StoredUserSettings>(SETTINGS_KEY, {})[user.id];
  return { ...user, ...settings };
}

export async function updateImageType(imageType: ImageType): Promise<void> {
  const settings = readStore<StoredUserSettings>(SETTINGS_KEY, {});
  writeStore(SETTINGS_KEY, { ...settings, [CURRENT_USER_ID]: { image_type: imageType } });
}
