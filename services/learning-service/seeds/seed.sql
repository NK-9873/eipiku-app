SET NAMES utf8mb4;

INSERT INTO courses (id, name, description) VALUES
  (1, '動物', '身近な動物の名前を覚えよう'),
  (2, '果物', 'よく見る果物の名前を覚えよう');

INSERT INTO words (id, english, japanese, example, photo_path, illustration_path) VALUES
  (1, 'dog', '犬', 'I have a dog.', '/images/photo/dog.jpg', '/images/illust/dog.png'),
  (2, 'cat', '猫', 'The cat is sleeping.', '/images/photo/cat.jpg', '/images/illust/cat.png'),
  (3, 'bird', '鳥', 'A bird is singing.', '/images/photo/bird.jpg', '/images/illust/bird.png'),
  (4, 'apple', 'りんご', 'I eat an apple every day.', '/images/photo/apple.jpg', '/images/illust/apple.png'),
  (5, 'banana', 'バナナ', 'Monkeys like bananas.', '/images/photo/banana.jpg', '/images/illust/banana.png'),
  (6, 'orange', 'オレンジ', 'This orange is sweet.', '/images/photo/orange.jpg', '/images/illust/orange.png');

INSERT INTO course_words (course_id, word_id, display_order) VALUES
  (1, 1, 1), (1, 2, 2), (1, 3, 3),
  (2, 4, 1), (2, 5, 2), (2, 6, 3);