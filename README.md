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

本番（EC2 1台 + Docker Compose）へのデプロイは [docs/deploy.md](docs/deploy.md) を参照。

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

## API 一覧

ログインが必要な API は、user-service の `POST /login` で取ったトークンを `Authorization: Bearer <token>` ヘッダーに付けて呼ぶ。
トークンの有効期限は24時間。どちらのサービスも同じ `JWT_SECRET` で署名・検証する。

権限の意味：**不要** = 誰でも呼べる / **ログイン** = トークンが必要（無い・不正・期限切れは 401）/ **管理者** = role が admin のユーザーのみ（それ以外は 403）

エラーは `{"error": "...", "details": ["..."]}` の JSON で返る（400 入力不正、401 未ログイン、403 権限なし、404 見つからない、409 重複）。

### user-service（http://localhost:8081）

| メソッド | URL | 内容 | 権限 |
|---|---|---|---|
| GET | /health | 動作確認 | 不要 |
| POST | /signup | 新規登録（name、email、password） | 不要 |
| POST | /login | ログイン。トークンを返す（email、password） | 不要 |
| GET | /me | 自分の情報 | ログイン |
| PUT | /me/settings | 画像タイプの変更（image_type: photo / illustration） | ログイン |
| GET | /users | ユーザー一覧 | 管理者 |
| GET | /users/{id} | ユーザーの詳細 | 管理者 |
| POST | /users | ユーザーの作成（name、email、password、role: admin / user。role 省略時は user） | 管理者 |
| PUT | /users/{id} | ユーザーの変更（name、email、role。password は入力したときだけ変更）。管理者が1人もいなくなる変更は 400 | 管理者 |

### learning-service（http://localhost:8082）

単語の情報（words・コースの単語・復習対象の単語）には、例文の日本語訳 `example_japanese` も含まれる。
単語の画像は、`POST /images` で返ってきた URL を `photo_path`・`illustration_path` に入れて保存する。
アップロードした画像は Docker のボリューム `learning_images` に保存される（`docker compose down -v` で消える）。

| メソッド | URL | 内容 | 権限 |
|---|---|---|---|
| GET | /health | 動作確認 | 不要 |
| GET | /courses | コース一覧 | ログイン |
| GET | /courses/{id} | コースの詳細 | ログイン |
| GET | /courses/{id}/words | コースの単語（出題順） | ログイン |
| POST | /courses | コースの登録 | 管理者 |
| PUT | /courses/{id} | コースの更新 | 管理者 |
| DELETE | /courses/{id} | コースの削除 | 管理者 |
| PUT | /courses/{id}/words | コースの単語と出題順をまとめて置き換え（word_ids の順が出題順） | 管理者 |
| GET | /words | 単語一覧（例文の日本語訳 example_japanese を含む） | ログイン |
| GET | /words/{id} | 単語の詳細 | ログイン |
| POST | /words | 単語の登録 | 管理者 |
| PUT | /words/{id} | 単語の更新 | 管理者 |
| DELETE | /words/{id} | 単語の削除 | 管理者 |
| POST | /images | 画像のアップロード（multipart/form-data の file 項目。JPEG・PNG・WebP、5MBまで）。表示用の URL を返す | 管理者 |
| GET | /images/{ファイル名} | アップロードした画像の表示（`<img>` から読み込むためトークン不要） | 不要 |
| GET | /me/word-progresses | 自分の習得度（?course_id=、?mastery=weak,unsure で絞り込み） | ログイン |
| PUT | /me/word-progresses/{wordId} | 習得度の登録・更新（mastery: weak / unsure / mastered） | ログイン |
| GET | /me/review-words | 復習対象の単語（?course_id=、?mastery=。mastery 省略時は苦手と微妙） | ログイン |
| GET | /me/course-progresses | 自分のコースの進み具合 | ログイン |
| PUT | /me/course-progresses/{courseId} | 中断位置の登録・更新（current_position） | ログイン |
| GET | /me/summary | マイページ用の学習状況 | ログイン |
| GET | /admin/users/{id}/progress | 指定ユーザーの習得度と進み具合 | 管理者 |

## フロントエンドの起動

```
cd frontend
npm install
npm run dev
```

http://localhost:3000 を開く。停止は `Ctrl+C`。

画面を使うには、先にバックエンドを起動しておく。開くとログイン画面になるので、新規登録するか、seed の管理者（admin@example.com）でログインする。
フロントは Next.js のサーバー側から API を呼ぶ。API の URL は `frontend/.env.local` の `USER_API_URL`・`LEARNING_API_URL` で指定する。
ログインのトークンは HttpOnly の Cookie（`eipiku_token`、24時間で期限切れ）に入れている。

## DB

### テーブル作成（migration）

```
docker run --rm -v "$(pwd)/services/user-service/migrations:/migrations" --network eipiku-app_default migrate/migrate -path=/migrations -database "mysql://users_app:<パスワード>@tcp(users-db:3306)/users" up
docker run --rm -v "$(pwd)/services/learning-service/migrations:/migrations" --network eipiku-app_default migrate/migrate -path=/migrations -database "mysql://learning_app:<パスワード>@tcp(learning-db:3306)/learning" up
```

`up` はまだ適用していない migration だけを順に適用する。migration が追加されたときも同じコマンドを実行すればよい。

| サービス | migration | 内容 |
|---|---|---|
| user-service | 000001 | users テーブルの作成 |
| learning-service | 000001 | courses・words・course_words・word_progresses・course_progresses テーブルの作成 |
| learning-service | 000002 | words に例文の日本語訳の列（example_japanese）を追加 |

**既に 000001 まで適用している DB の場合**：上の learning-service のコマンドを実行すると 000002 だけが適用される。
そのあと次の「初期データ投入」の learning-service のコマンドを実行すると、seed の6語に例文の日本語訳が入る
（learning-service の seed は何度実行してもよい。既にあるコース・単語は変えずに、例文の日本語訳だけを入れる）。

今どこまで適用されているかは `up` の代わりに `version` で確認できる。

### 初期データ投入（seed）

```
docker compose exec -T learning-db sh -c 'mysql -u learning_app -p"$MYSQL_PASSWORD" learning' < services/learning-service/seeds/seed.sql
docker compose exec -T users-db sh -c 'mysql -u users_app -p"$MYSQL_PASSWORD" users' < services/user-service/seeds/seed.sql
```

user-service の seed では管理者ユーザーが1人登録される（こちらは2回目に実行すると重複エラーになる）。

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