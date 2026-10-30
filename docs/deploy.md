# デプロイ手順（EC2 1台 + Docker Compose）

えいピクを AWS の EC2 1台に Docker Compose で動かす手順。
当面は独自ドメインと HTTPS を使わず、`http://<EC2 のパブリック IP>`（80番）で開く。

## 構成

```
ブラウザ ──http:80──▶ frontend（Next.js）
                         │  compose の内部ネットワーク（サービス名で通信）
                         ├──▶ user-service ──▶ users-db
                         └──▶ learning-service ──▶ learning-db
                                    └ 画像（ボリューム learning_images）
```

- インターネットに公開するのは frontend の 80番だけ。API と DB はポートを公開しない（BFF 構成）。
- 単語の画像は `http://<IP>/images/<ファイル名>` で開く。frontend が learning-service に転送する（`frontend/app/images/[name]/route.ts`）。
- DB のデータと画像は Docker のボリューム（`eipiku_users_data`・`eipiku_learning_data`・`eipiku_learning_images`）に保存される。コンテナを作り直しても消えない。

使うファイル：

| ファイル | 内容 |
|---|---|
| `compose.prod.yaml` | 本番用の Compose ファイル |
| `.env.prod.example` | 本番用の環境変数の見本。EC2 上で `.env.prod` にコピーして使う |
| `frontend/Dockerfile` | frontend の本番用イメージ（`next build` して本番モードで起動） |

## 1. EC2 の準備（AWS コンソール）

- **AMI**：Amazon Linux 2023（リージョン ap-southeast-2）
- **インスタンスタイプ**：t3.small（メモリ 2GB）以上が目安。t3.micro（1GB）でも動くが、ビルドにはスワップが必須（手順 3）
- **ストレージ**：20GB 以上（標準の 8GB だと、イメージとビルドのキャッシュですぐ一杯になる）
- **セキュリティグループ**（受信ルール）：

  | ポート | 送信元 | 用途 |
  |---|---|---|
  | 22（SSH） | 自分の IP だけ | ログイン |
  | 80（HTTP） | 0.0.0.0/0 | アプリ |

  3306・8080 などは開けない。
- **Elastic IP**：割り当てておくのがおすすめ。割り当てないと、インスタンスを停止・起動するたびにパブリック IP が変わり、開く URL も変わる。

SSH でログインする：

```
ssh -i <キーファイル>.pem ec2-user@<パブリック IP>
```

以下はすべて EC2 上で実行する。

## 2. Docker と Docker Compose のインストール

```
sudo dnf update -y
sudo dnf install -y docker git
sudo systemctl enable --now docker
sudo usermod -aG docker ec2-user   # sudo なしで docker を使えるようにする
```

Amazon Linux 2023 の docker パッケージには Compose と Buildx（Compose のビルドに必要）が入っていないので、別に入れる。

```
sudo mkdir -p /usr/local/lib/docker/cli-plugins

# Docker Compose
sudo curl -fL -o /usr/local/lib/docker/cli-plugins/docker-compose \
  "https://github.com/docker/compose/releases/latest/download/docker-compose-linux-$(uname -m)"

# Docker Buildx（ファイル名にバージョンと amd64 / arm64 が入るので、最新版の番号を調べてから取る）
ARCH=$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
BUILDX_VERSION=$(curl -fsSL https://api.github.com/repos/docker/buildx/releases/latest | grep -m1 '"tag_name"' | cut -d '"' -f4)
sudo curl -fL -o /usr/local/lib/docker/cli-plugins/docker-buildx \
  "https://github.com/docker/buildx/releases/download/${BUILDX_VERSION}/buildx-${BUILDX_VERSION}.linux-${ARCH}"

sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose /usr/local/lib/docker/cli-plugins/docker-buildx
```

**一度ログアウトして SSH で入り直す**（docker グループの設定を反映するため）。入り直したら確認する：

```
docker version
docker compose version
docker buildx version
```

## 3. スワップの作成（メモリが 2GB 以下のとき）

frontend のビルド（`next build`）はメモリを多く使い、メモリが足りないと途中で止まる（`exit code: 137` や `Killed`）。
メモリが少ない EC2 では、ディスクの一部をメモリの代わりに使うスワップを 2GB 作っておく。

```
sudo dd if=/dev/zero of=/swapfile bs=1M count=2048
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab   # 再起動後も使う
free -h   # Swap の行が 2.0Gi になっていれば OK
```

## 4. リポジトリの取得

```
git clone https://github.com/NK-9873/eipiku-app.git
cd eipiku-app
```

デプロイするブランチを使う（通常は main）。以降のコマンドはすべて `eipiku-app` ディレクトリで実行する。

## 5. .env.prod の作成

```
cp .env.prod.example .env.prod
chmod 600 .env.prod   # 自分だけが読めるようにする
```

パスワードと `JWT_SECRET` には、**必ず推測されない強いランダム値**を入れる。
次のコマンドで `change_me` をすべてランダムな値に置き換え、`PUBLIC_ORIGIN` を EC2 のパブリック IP にする。

```
PUBLIC_IP=<EC2 のパブリック IP（Elastic IP）>
sed -i \
  -e "s|^PUBLIC_ORIGIN=.*|PUBLIC_ORIGIN=http://${PUBLIC_IP}|" \
  -e "s|^JWT_SECRET=.*|JWT_SECRET=$(openssl rand -base64 48 | tr -d '\n')|" \
  -e "s|^USERS_DB_PASSWORD=.*|USERS_DB_PASSWORD=$(openssl rand -hex 24)|" \
  -e "s|^USERS_ROOT_PASSWORD=.*|USERS_ROOT_PASSWORD=$(openssl rand -hex 24)|" \
  -e "s|^LEARNING_DB_PASSWORD=.*|LEARNING_DB_PASSWORD=$(openssl rand -hex 24)|" \
  -e "s|^LEARNING_ROOT_PASSWORD=.*|LEARNING_ROOT_PASSWORD=$(openssl rand -hex 24)|" \
  .env.prod
grep -c change_me .env.prod   # 0 なら置き換え済み
```

| 変数 | 内容 |
|---|---|
| `PUBLIC_ORIGIN` | ブラウザから開く URL（`http://<IP>`） |
| `HTTP_PORT` | frontend を公開するポート（80） |
| `COOKIE_SECURE` | ログインの Cookie に Secure 属性を付けるか。**HTTPS なしの間は `false`**（`true` だとログインできない） |
| `IMAGE_BASE_URL` | アップロードした画像の URL の先頭。通常は `/images` のまま |
| `JWT_SECRET` | ログイン用トークンの署名の鍵（32文字以上） |
| `USERS_DB_PASSWORD` / `USERS_ROOT_PASSWORD` | users-db のパスワード |
| `LEARNING_DB_PASSWORD` / `LEARNING_ROOT_PASSWORD` | learning-db のパスワード |

注意：

- `.env.prod` は Git に入れない（`.gitignore` 済み）。中身をチャットやリポジトリに貼らない。
- DB のパスワードは **最初に DB を作ったときにだけ** 設定される。あとから `.env.prod` を変えても DB 側は変わらず、API が DB に接続できなくなる（変えるときは DB の中でもパスワードを変更する）。
- DB のパスワードは migration の接続 URL にも入るので、記号を含まない hex（`openssl rand -hex`）にしている。

## 6. 起動

コマンドが長いので、エイリアスを作っておく（以降は `dc` を使う）。

```
echo "alias dc='docker compose -f compose.prod.yaml --env-file .env.prod'" >> ~/.bashrc
source ~/.bashrc
```

ビルドして起動する（初回は数分〜十数分かかる）：

```
dc up -d --build
dc ps   # 5つとも Up、DB は (healthy) になるまで待つ
```

## 7. テーブル作成（migration）と初期データ（seed）

テーブルを作る（初回と、migration が追加されたとき）：

```
dc run --rm migrate-users
dc run --rm migrate-learning
```

まだ適用していない migration だけが適用される。何度実行してもよい（新しいものが無ければ `no change`）。

サンプルのコース・単語を入れる（任意）：

```
dc exec -T learning-db sh -c 'mysql -u learning_app -p"$MYSQL_PASSWORD" learning' < services/learning-service/seeds/seed.sql
```

seed の6語は画像が「画像なし」になる（仮のパスのため）。管理画面から画像を登録する。

### 管理者ユーザーの作成

**user-service の seed（`services/user-service/seeds/seed.sql`）は本番では実行しない。**
管理者のパスワードがリポジトリに公開されているため。代わりに次の手順で管理者を作る。

1. ブラウザで `http://<IP>/signup` を開き、自分のアカウントを新規登録する（パスワードは強いものにする）
2. EC2 上で、そのアカウントを管理者にする：

   ```
   dc exec users-db sh -c 'mysql -u users_app -p"$MYSQL_PASSWORD" users -e "UPDATE users SET role = '"'"'admin'"'"' WHERE email = '"'"'<登録したメールアドレス>'"'"'; SELECT id, email, role FROM users;"'
   ```

   一覧で role が `admin` になっていれば OK。
3. ログインし直すと、ヘッダーに「管理画面」が出る。2人目以降の管理者は管理画面の「ユーザー」から作れる。

## 8. 動作確認

ブラウザで `http://<IP>` を開き、次を確かめる：

- ログイン画面が出て、ログインできる
- コースを開いて学習できる
- 管理画面で単語に画像を登録でき、一覧と学習画面に画像が出る

うまくいかないときはログを見る：

```
dc ps
dc logs --tail=50 frontend
dc logs --tail=50 learning-service
```

## 9. 更新（新しいコードを反映する）

```
cd ~/eipiku-app
git pull
dc up -d --build                # 変更のあったサービスだけ作り直す
dc run --rm migrate-users       # migration が追加されていれば適用される
dc run --rm migrate-learning
docker image prune -f           # 古いイメージを消してディスクを空ける
```

DB と画像はボリュームにあるので、更新しても消えない。

## よく使うコマンド

```
dc ps                         # 状態
dc logs -f frontend           # ログを見続ける（Ctrl+C で抜ける）
dc restart learning-service   # 再起動
dc down                       # 停止（データは残る）
dc up -d                      # 起動
df -h /                       # ディスクの空き
docker system df              # Docker が使っているディスク
```

**`dc down -v` は DB と画像をすべて消す。** 本番では使わない。

## バックアップ（任意）

DB のダンプと画像を `~/backup` に保存する：

```
mkdir -p ~/backup
dc exec -T users-db sh -c 'mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" --single-transaction users' > ~/backup/users-$(date +%F).sql
dc exec -T learning-db sh -c 'mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" --single-transaction learning' > ~/backup/learning-$(date +%F).sql
dc exec -T learning-service tar czf - -C /data images > ~/backup/images-$(date +%F).tar.gz
```

EC2 が壊れると一緒に消えるので、大事なデータは S3 や手元の PC にもコピーしておく。EBS のスナップショットでもよい。

## よくあるエラー

- **ビルドが `exit code: 137` や `Killed` で止まる**：メモリ不足。手順 3 のスワップを作る
- **`compose build requires buildx` と出る**：手順 2 の Buildx が入っていない
- **ログインしてもすぐログイン画面に戻る**：`.env.prod` の `COOKIE_SECURE` が `true` になっている。HTTPS なしの間は `false` にして `dc up -d` で作り直す
- **`.env.prod に設定してください` と出る**：`.env.prod` にその変数が無い・空。`--env-file .env.prod` を付け忘れていないかも確認する（`dc` エイリアスを使う）
- **API が DB に接続できない（`Access denied`）**：DB を作ったあとで `.env.prod` の DB のパスワードを変えた。手順 5 の注意を参照
- **ブラウザで開けない**：セキュリティグループで 80番が開いているか、URL が `https://` ではなく `http://` になっているかを確認する
- **ディスクが一杯**：`docker image prune -f` と `docker builder prune -f` で古いイメージとビルドのキャッシュを消す

## 今後 HTTPS にするとき

独自ドメインと HTTPS（ALB + ACM、または Caddy・nginx など）にしたら、`.env.prod` を次のように変えて `dc up -d` で作り直す。

```
PUBLIC_ORIGIN=https://<ドメイン>
COOKIE_SECURE=true
```
