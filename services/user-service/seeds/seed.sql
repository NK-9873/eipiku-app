SET NAMES utf8mb4;

-- ローカル開発用の管理者ユーザー（本番では使わないこと）
-- メールアドレス: admin@example.com / パスワード: eipiku-admin-dev
-- password_hash はこのパスワードを bcrypt（コスト10）でハッシュ化したもの
INSERT INTO users (id, name, email, password_hash, image_type, role) VALUES
  (1, '管理者', 'admin@example.com', '$2a$10$FdLnQdBBWuEzeXPPDsYN2ON2GDyQRuAyded1Wt4OPJYW5a/5Ao41q', 'photo', 'admin');
