// services/learning-service/seeds/seed.sql と同じ内容のモックデータ。
// API ができるまでの仮データなので、カラム名もテーブル定義に合わせている。

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
  photo_path: string | null;
  illustration_path: string | null;
};

export type CourseWord = {
  course_id: number;
  word_id: number;
  display_order: number;
};

export const courses: Course[] = [
  { id: 1, name: "動物", description: "身近な動物の名前を覚えよう" },
  { id: 2, name: "果物", description: "よく見る果物の名前を覚えよう" },
];

export const words: Word[] = [
  { id: 1, english: "dog", japanese: "犬", example: "I have a dog.", photo_path: "/images/photo/dog.jpg", illustration_path: "/images/illust/dog.png" },
  { id: 2, english: "cat", japanese: "猫", example: "The cat is sleeping.", photo_path: "/images/photo/cat.jpg", illustration_path: "/images/illust/cat.png" },
  { id: 3, english: "bird", japanese: "鳥", example: "A bird is singing.", photo_path: "/images/photo/bird.jpg", illustration_path: "/images/illust/bird.png" },
  { id: 4, english: "apple", japanese: "りんご", example: "I eat an apple every day.", photo_path: "/images/photo/apple.jpg", illustration_path: "/images/illust/apple.png" },
  { id: 5, english: "banana", japanese: "バナナ", example: "Monkeys like bananas.", photo_path: "/images/photo/banana.jpg", illustration_path: "/images/illust/banana.png" },
  { id: 6, english: "orange", japanese: "オレンジ", example: "This orange is sweet.", photo_path: "/images/photo/orange.jpg", illustration_path: "/images/illust/orange.png" },
];

export const course_words: CourseWord[] = [
  { course_id: 1, word_id: 1, display_order: 1 },
  { course_id: 1, word_id: 2, display_order: 2 },
  { course_id: 1, word_id: 3, display_order: 3 },
  { course_id: 2, word_id: 4, display_order: 1 },
  { course_id: 2, word_id: 5, display_order: 2 },
  { course_id: 2, word_id: 6, display_order: 3 },
];
