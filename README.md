# URL Shortener API (Go + Echo + Redis)

短縮URLを生成・リダイレクトするAPIサービス。  
Go (Echo) + Redis、Docker Compose でマルチコンテナ構成を採用。  
Swagger UIでAPI仕様をブラウザから確認可能。

---

## 特徴

- `/shorten` で長いURLを短縮（オプションでカスタムコードも指定可）
- `/r/{code}` でリダイレクト
- Redisキャッシュによる高速化
- `/healthz` での健康チェック
- OpenAPI（Swagger UI）対応
- HEADメソッドでもリダイレクト先の `Location` を返却
- 簡易レートリミット（60 req/min）
- ユニットテスト例を同梱

---

## 起動方法

```bash
docker compose up --build
起動後:

健康チェック: http://localhost:8080/healthz → ok

Swagger UI: http://localhost:8080/docs

Swagger JSON: http://localhost:8080/swagger.json

使用例（cmd.exe）
REM 短縮URLを作成
curl -X POST http://localhost:8080/shorten ^
  -H "Content-Type: application/json" ^
  -d "{\"url\": \"https://example.com/some/very/long/path\"}"

REM 短縮コードを使ってアクセス
curl -i http://localhost:8080/r/99jnEc

REM HEADでリダイレクト先だけ確認
curl -IL http://localhost:8080/r/99jnEc
Redisでの確認
docker exec -it urlshort-redis redis-cli
keys *
get short:<code>
短縮コード → 元URL を Redis にキャッシュ（永続ストレージにも保存可）

ステータスコード方針:

短縮作成成功: 201 Created

リダイレクト成功: 301 Moved Permanently

存在しないコード: 404 Not Found

バリデーションエラー: 400 Bad Request

Swagger UI
ブラウザで http://localhost:8080/docs

画面上でエンドポイントの動作を確認可能

ユニットテスト（例）
internal/core/service_test.go に Shorten / Resolve のテーブルテスト例を収録。

実行:

go test ./...

フォルダ構成
urlshort/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── app/server.go
│   ├── core/service.go
│   ├── core/service_test.go
│   └── infra/redis.go
├── docs/swagger.json
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md

---

## GitHub Push 手順（cmd.exe）

```bat
REM 1. 作業ディレクトリへ移動
cd C:\urlshort

REM 2. 最新差分を追加
git add .

REM 3. コミット（メッセージは適宜変更）
git commit -m "feat: add Swagger UI, health check, HEAD method support"

REM 4. 最新のmainブランチを取り込み（コンフリクト回避）
git pull --rebase origin main

REM 5. リモートへPush
git push origin main

⚠️ git pull --rebase でエラーが出た場合は、競合解消後に再度 git rebase --continue を実行してから git push します。