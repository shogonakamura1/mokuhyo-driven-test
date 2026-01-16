# 目標ドリブン・フローティングノート

目標に向かって思考を深掘りし、具体的な行動に落とし込むためのアプリケーションです。

## 技術スタック

### フロントエンド
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- Supabase Auth

### バックエンド
- Go 1.21+
- Gin (Web Framework)
- PostgreSQL (Supabase)
- pgx (PostgreSQL driver)

### インフラ
- Supabase (認証・データベース)
- Cloudflare Pages / Vercel (フロントエンド)
- Fly.io / Render / Railway (バックエンド)

## プロジェクト構成

```
goal-tree/
├── apps/
│   ├── web/              # Next.jsフロントエンド
│   └── api/              # Goバックエンド
├── supabase/
│   └── migrations/       # データベースマイグレーション
└── packages/
    └── shared/           # 共有型定義（将来拡張用）
```

## セットアップ

### 1. 前提条件

- Node.js 18+
- Go 1.21+
- Supabaseアカウント

### 2. Supabaseの設定

1. [Supabase](https://supabase.com)でプロジェクトを作成
2. プロジェクトのURLとAnon Keyを取得
3. Database > SQL Editorでマイグレーションファイルを実行:
   ```bash
   cat supabase/migrations/001_initial_schema.sql
   ```
4. Authentication > ProvidersでGoogle認証を有効化

### 3. バックエンド（Go API）のセットアップ

```bash
cd apps/api

# 依存関係のインストール
go mod tidy
go mod download

# 環境変数の設定
cp .env.example .env
# .envファイルを編集して以下を設定:
# - DATABASE_URL: SupabaseのPostgreSQL接続文字列
# - SUPABASE_JWKS_URL: https://<project-ref>.supabase.co/auth/v1/keys
# - PORT: 8080 (デフォルト)

# サーバーの起動
go run cmd/server/main.go
```

### 4. フロントエンド（Next.js）のセットアップ

```bash
cd apps/web

# 依存関係のインストール
npm install

# 環境変数の設定
# .env.localファイルを作成して以下を設定:
# NEXT_PUBLIC_SUPABASE_URL=https://<project-ref>.supabase.co
# NEXT_PUBLIC_SUPABASE_ANON_KEY=<anon-key>
# NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# 開発サーバーの起動
npm run dev
```

## 環境変数

### バックエンド（apps/api/.env）

```env
DATABASE_URL=postgresql://user:password@host:5432/dbname
SUPABASE_JWKS_URL=https://your-project.supabase.co/auth/v1/keys
PORT=8080
```

### フロントエンド（apps/web/.env.local）

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## 使い方

### 1. ログイン
- `/login` にアクセスしてGoogleアカウントでログイン

### 2. プロジェクト作成
- `/input` で目標のタイトルを入力してプロジェクトを作成

### 3. ツリー編集
- `/p/[projectId]` でツリーを編集
- **Enter**: 編集確定 → 兄弟ノード追加
- **Tab**: 編集確定 → 子ノード追加
- **Shift+Enter**: 編集確定のみ
- **Esc**: 編集キャンセル
- **右クリック**: ノード削除（子孫も含む）

## APIエンドポイント

### 認証
- `GET /v1/me` - 現在のユーザー情報と設定を取得

### プロジェクト
- `POST /v1/projects` - プロジェクト作成
- `GET /v1/projects` - プロジェクト一覧取得
- `GET /v1/projects/:projectId` - プロジェクト詳細取得
- `PATCH /v1/projects/:projectId` - プロジェクト更新
- `GET /v1/projects/:projectId/tree` - ツリー構造取得
- `POST /v1/projects/:projectId/save` - 保存マーク

### ノード
- `POST /v1/projects/:projectId/nodes` - ノード作成
- `PATCH /v1/projects/:projectId/nodes/:nodeId` - ノード更新
- `DELETE /v1/projects/:projectId/nodes/:nodeId` - ノード削除（論理削除、子孫含む）

### エッジ
- `PATCH /v1/projects/:projectId/edges/:edgeId` - エッジ更新（関係ラベル）
- `POST /v1/projects/:projectId/reorder` - ノードの並び替え

### 設定
- `GET /v1/settings` - ユーザー設定取得
- `PATCH /v1/settings` - ユーザー設定更新

## デプロイ

### フロントエンド（Cloudflare Pages / Vercel）

1. リポジトリを接続
2. ビルドコマンド: `cd apps/web && npm run build`
3. 出力ディレクトリ: `apps/web/.next`
4. 環境変数を設定

### バックエンド（Fly.io）

```bash
cd apps/api
fly launch
# 環境変数を設定
fly secrets set DATABASE_URL=...
fly secrets set SUPABASE_JWKS_URL=...
```

## 開発

### データベースマイグレーション

SupabaseのSQL Editorでマイグレーションファイルを実行してください。

### テスト

```bash
# バックエンド
cd apps/api
go test ./...

# フロントエンド
cd apps/web
npm test
```

## ライセンス

MIT
