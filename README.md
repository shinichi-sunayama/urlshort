# 🔗 urlshort - シンプルなURL短縮サービス

GoとRedisを使って構築した、シンプルかつ拡張性の高い**URL短縮サービス**です。  
フォームUI、カスタムコード対応、キャッシュ対応、Docker構成などを実装済み。

---

## 🚀 主な特徴

- ✅ 任意のURLを短縮し、リダイレクト可能なコードに変換
- ✅ Redisによる高速な永続保存＆キャッシュ対応
- ✅ カスタムURLコードの指定も可能（例：`/r/my-code`）
- ✅ BootstrapでスタイリングされたUI
- ✅ コピー用ボタン付き（短縮URLをワンクリックでコピー）
- ✅ Docker & docker-compose 対応
- ✅ Web UIはGoの `html/template` で構築
- ✅ SwaggerでAPIドキュメント提供可能（内部にjsonあり）

---

## 🧱 使用技術スタック

| 分類          | 技術                                      |
|---------------|-------------------------------------------|
| 言語           | Go 1.20+                                  |
| Webフレーム    | [Echo](https://echo.labstack.com/) v4     |
| テンプレート   | `html/template`                           |
| ストレージ     | Redis（キャッシュ＋永続HSET）            |
| キャッシュTTL  | 環境変数で指定（デフォルト：3600秒）      |
| API管理        | Swagger JSON（`internal/app/docs`）       |
| スタイル       | Bootstrap 5.x                             |

---

## 📁 ディレクトリ構成

```plaintext
urlshort/
├── cmd/server/
│   ├── main.go             # 起動エントリ
│   ├── .env                # 接続先など設定
│   └── templates/
│       └── index.html      # UIテンプレート（フォーム＋結果表示）
│
├── internal/
│   ├── app/
│   │   ├── server.go       # ルーティング・UIハンドラ
│   │   ├── renderer.go     # Echoテンプレートレンダラ
│   │   └── docs/
│   │       └── swagger.json# APIドキュメント
│   │
│   ├── logic/
│   │   ├── service.go      # コアビジネスロジック
│   │   ├── codegen.go      # ランダムコード生成
│   │   └── service_test.go # 単体テスト
│   │
│   └── store/
│       └── redis_store.go  # Redisとのやりとり
│
├── docker-compose.yml      # Redis含むサービス構成
├── Dockerfile              # Goアプリのビルド用
├── go.mod / go.sum         # モジュール定義
├── LICENSE
├── server.exe              # Windowsビルド済みバイナリ（例）
└── README.md               # ← 本ファイル
🖥 画面イメージ
✔️ URLフォームに入力 → 短縮ボタン → 結果表示

✔️ 短縮URLはクリックでコピー可能

✔️ Bootstrapでスタイリング済み

🔧 .env構成
以下のような .env を cmd/server/ に配置してください：

PORT=8081
BASE_URL=http://localhost:8081
REDIS_ADDR=localhost:6379
REDIS_DB=0
CACHE_TTL_SECONDS=3600
REDIS_URL=redis:6379
📦 起動方法（ローカル）
✅ 前提
Redisを localhost:6379 または .envの設定に準拠して起動済みであること

✅ 起動
cd cmd/server
go run main.go
→ http://localhost:8081 にアクセス

🐳 Docker / docker-compose

Redis付きで立ち上げるには：
docker-compose up --build

またはポート競合時：
services:
  redis:
    ports:
      - "6380:6379"  # ホスト6380 → コンテナ6379
.env 側も合わせて REDIS_ADDR=localhost:6380 にしてください。

✅ 主なルーティング
Method	Path	説明
GET	/	入力フォーム画面
POST	/	URLを短縮（フォーム送信）
GET	/r/:code	短縮コードからリダイレクト

✨ 実装済み機能一覧
 入力フォーム

 短縮コード自動生成

 任意コードの指定対応

 Redis保存

 Redisキャッシュ

 URLのコピーボタン（JS）

 エラーハンドリング

 Bootstrap対応のUI

 Swagger API JSONファイル（OpenAPI 3）