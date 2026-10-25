# えいピク（eipiku-app）

写真かイラストを選んで覚える英単語帳アプリ

## 構成

| サービス | 役割 | PC側のポート |
|---|---|---|
| frontend | 画面（Next.js） | 3000 |
| user-service | ユーザー登録・ログイン・設定 | 8081 |
| learning-service | コース・単語学習・復習・管理 | 8082 |
| users-db | user-service 用 MySQL | 13306 |
| learning-db | learning-service 用 MySQL | 13307 |

## 必要なツール

- Docker Desktop（Docker Compose を含む）
- Git
- Node.js（frontend を動かす場合）
- Go 1.25（コンテナ外で Go を動かす場合のみ）

## 初期設定

```
git clone https://github.com/NK-9873/eipiku-app.git
cd eipiku-app
cp .env.example .env
cp services/user-service/.env.example services/user-service/.env
cp services/learning-service/.env.example services/learning-service/.env
cp frontend/.env.local.example frontend/.env.local
```

各 .env の `change_me` をローカル用のパスワードに書き換える。
サービス側の `DB_PASSWORD` は、リポジトリ直下の .env の対応するパスワードと同じ値にする。

`JWT_SECRET`（ログイン用トークンの署名の鍵）は32文字以上のランダムな文字列にし、
user-service と learning-service の .env で**同じ値**にする。次のコマンドで作れる。

```
openssl rand -base64 48
```

## バックエンドの起動・停止

```
docker compose up --build -d    # 起動（コード変更の反映もこれ）
docker compose ps               # 状態確認
docker compose logs --tail=50 user-service   # ログ確認
docker compose down             # 停止（データは残る）
```

`docker compose down -v` は DB のデータも削除するため、初期化したいとき以外は使わない。

動作確認：

- http://localhost:8081/health
- http://localhost:8082/health
- http://localhost:8082/courses（コース一覧の JSON）

## フロントエンドの起動

```
cd frontend
npm install
npm run dev
```

http://localhost:3000 を開く。停止は `Ctrl+C`。

## DB

### テーブル作成（migration）

```
docker run --rm -v "$(pwd)/services/user-service/migrations:/migrations" --network eipiku-app_default migrate/migrate -path=/migrations -database "mysql://users_app:<パスワード>@tcp(users-db:3306)/users" up
docker run --rm -v "$(pwd)/services/learning-service/migrations:/migrations" --network eipiku-app_default migrate/migrate -path=/migrations -database "mysql://learning_app:<パスワード>@tcp(learning-db:3306)/learning" up
```

### 初期データ投入（seed）

```
docker compose exec -T learning-db sh -c 'mysql -u learning_app -p"$MYSQL_PASSWORD" learning' < services/learning-service/seeds/seed.sql
docker compose exec -T users-db sh -c 'mysql -u users_app -p"$MYSQL_PASSWORD" users' < services/user-service/seeds/seed.sql
```

user-service の seed では管理者ユーザーが1人登録される。

| メールアドレス | パスワード | role |
|---|---|---|
| admin@example.com | eipiku-admin-dev | admin |

**このパスワードはローカル開発用。** リポジトリに公開されているので、本番などローカル以外の環境では使わないこと。

### DB に接続して中身を見る

```
docker compose exec learning-db mysql -u learning_app -p --default-character-set=utf8mb4 learning
```

## よくあるエラー

- **起動直後に API が動かない**：DB の準備待ち。`docker compose ps` で DB が healthy になるまで待つ
- **日本語が ?? になる**：接続時に `--default-character-set=utf8mb4` をつける
- **.env を変えたのに反映されない**：`docker compose restart` ではなく `docker compose up -d` で作り直す
- **ポートが使用中**：他のアプリが 3000、8081、13306 などを使っていないか確認する