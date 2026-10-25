// API（user-service・learning-service）が返すデータの型と、画面で使う表示名。
// サーバー・ブラウザのどちらからも読み込んでよい（API の呼び出し自体は lib/api.ts などサーバー側だけで行う）

export type Mastery = "weak" | "unsure" | "mastered";
export type ImageType = "photo" | "illustration";
export type Role = "user" | "admin";

export const MASTERY_LABELS: Record<Mastery, string> = {
  weak: "苦手",
  unsure: "微妙",
  mastered: "覚えた",
};

export const IMAGE_TYPE_LABELS: Record<ImageType, string> = {
  photo: "写真",
  illustration: "イラスト",
};

export type ReviewTarget = "weak" | "unsure" | "both";

export const REVIEW_TARGETS: { value: ReviewTarget; label: string; masteries: Mastery[] }[] = [
  { value: "weak", label: "苦手", masteries: ["weak"] },
  { value: "unsure", label: "微妙", masteries: ["unsure"] },
  { value: "both", label: "両方", masteries: ["weak", "unsure"] },
];

// user-service
export type User = {
  id: number;
  name: string;
  email: string;
  image_type: ImageType;
  role: Role;
  created_at: string;
  updated_at: string;
};

// learning-service
export type Course = {
  id: number;
  name: string;
  description: string | null;
};

export type Word = {
  id: number;
  english: string;
  japanese: string;
  example: string | null;
  example_japanese: string | null;
  photo_path: string | null; // 画像の URL（learning-service の POST /images が返したもの）
  illustration_path: string | null;
};

// GET /courses/{id}/words の1件
export type CourseWord = Word & { display_order: number };

// GET /me/review-words の1件
export type ReviewWord = Word & { mastery: Mastery; updated_at: string };

// current_position は「何語目まで学習したか」（0 なら未学習、単語数と同じなら最後まで終わっている）
export type CourseProgress = {
  user_id: number;
  course_id: number;
  current_position: number;
  updated_at: string;
};

export type CourseStatus = "not_started" | "in_progress" | "completed";

// GET /me/summary
export type Summary = {
  learned_words: number;
  total_words: number;
  weak: number;
  unsure: number;
  mastered: number;
  courses: {
    course_id: number;
    name: string;
    total_words: number;
    current_position: number;
    status: CourseStatus;
    updated_at: string | null;
  }[];
};

// API のエラー（{"error": "...", "details": [...]}）を画面に出すための形
export type FormState = {
  error?: string;
  details?: string[];
  // エラーのとき入力し直さずに済むよう、送信された値（パスワード以外）を返す
  values?: { name?: string; email?: string };
};
