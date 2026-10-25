SET NAMES utf8mb4;

-- 何度実行してもよいようにしている。
-- コースとコースの単語は、既にあれば何もしない。単語は、既にあれば例文の日本語訳だけを入れる
-- （000002 の migration で列を追加した既存の DB にも、訳を入れられるようにするため）

INSERT IGNORE INTO courses (id, name, description) VALUES
  (1, '動物', '身近な動物の名前を覚えよう'),
  (2, '果物', 'よく見る果物の名前を覚えよう');

INSERT INTO words (id, english, japanese, example, example_japanese, photo_path, illustration_path) VALUES
  (1, 'dog', '犬', 'I have a dog.', '私は犬を飼っています。', '/images/photo/dog.jpg', '/images/illust/dog.png'),
  (2, 'cat', '猫', 'The cat is sleeping.', '猫が眠っています。', '/images/photo/cat.jpg', '/images/illust/cat.png'),
  (3, 'bird', '鳥', 'A bird is singing.', '鳥がさえずっています。', '/images/photo/bird.jpg', '/images/illust/bird.png'),
  (4, 'apple', 'りんご', 'I eat an apple every day.', '私は毎日りんごを食べます。', '/images/photo/apple.jpg', '/images/illust/apple.png'),
  (5, 'banana', 'バナナ', 'Monkeys like bananas.', 'サルはバナナが好きです。', '/images/photo/banana.jpg', '/images/illust/banana.png'),
  (6, 'orange', 'オレンジ', 'This orange is sweet.', 'このオレンジは甘いです。', '/images/photo/orange.jpg', '/images/illust/orange.png')
AS new
ON DUPLICATE KEY UPDATE example_japanese = new.example_japanese;

INSERT IGNORE INTO course_words (course_id, word_id, display_order) VALUES
  (1, 1, 1), (1, 2, 2), (1, 3, 3),
  (2, 4, 1), (2, 5, 2), (2, 6, 3);
