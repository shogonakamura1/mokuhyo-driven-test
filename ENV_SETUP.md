# 環境変数設定ガイド

このプロジェクトでは、バックエンド（Go API）とフロントエンド（Next.js）でそれぞれ環境変数ファイルが必要です。

## 1. Supabaseの設定値を取得する

### 1.1 Supabaseプロジェクトを作成
1. [Supabase](https://supabase.com)にアクセスしてアカウントを作成
2. 新しいプロジェクトを作成
3. プロジェクトの設定を開く

### 1.2 必要な値を取得

#### Project URL
- **場所**: Settings > API > Project URL
- **例**: `https://abcdefghijklmnop.supabase.co`
- **用途**: フロントエンドの認証とAPI接続

#### Anon Key（公開キー）
- **場所**: Settings > API > Project API keys > `anon` `public`
- **例**: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- **用途**: フロントエンドのSupabaseクライアント初期化

#### Database Connection String
- **場所**: Settings > Database > Connection string > URI
- **例**: `postgresql://postgres:[YOUR-PASSWORD]@db.abcdefghijklmnop.supabase.co:5432/postgres`
- **用途**: バックエンドのデータベース接続
- **注意**: `[YOUR-PASSWORD]` の部分を実際のパスワードに置き換える必要があります

#### Project Reference
- **場所**: Settings > General > Reference ID
- **例**: `abcdefghijklmnop`
- **用途**: JWKS URLの構築

## 2. バックエンド（Go API）の環境変数

### ファイル場所
```
apps/api/.env
```

### 設定内容

#### ローカル開発時（ローカルPostgreSQLを使用）

```env
# データベースタイプ
# local または postgres: ローカルPostgreSQLを使用（開発時）
# supabase: Supabaseを使用（デプロイ時）
# デフォルト: local（省略可能）
DB_TYPE=local

# データベース接続文字列（ローカルPostgreSQL）
# 形式: postgresql://user:password@localhost:5432/dbname
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mokuhyo_db_test

# Supabase JWKS URL（JWT検証用）
# 認証はSupabaseを使用するため、ローカル開発時でも必要
SUPABASE_JWKS_URL=https://your-project-ref.supabase.co/auth/v1/keys

# APIサーバーのポート番号（デフォルト: 8080）
PORT=8080
```

#### デプロイ時（Supabaseを使用）

```env
# データベースタイプ
DB_TYPE=supabase

# データベース接続文字列（Supabase PostgreSQL）
# 形式: postgresql://postgres:password@db.project-ref.supabase.co:5432/postgres
DATABASE_URL=postgresql://postgres:your-password@db.your-project-ref.supabase.co:5432/postgres

# Supabase JWKS URL（JWT検証用）
SUPABASE_JWKS_URL=https://your-project-ref.supabase.co/auth/v1/keys

# APIサーバーのポート番号（デフォルト: 8080）
PORT=8080
```

### 設定手順

1. `apps/api` ディレクトリに `.env` ファイルを作成
   ```bash
   cd apps/api
   touch .env
   ```

2. `.env` ファイルを編集して実際の値を設定

   **ローカル開発時:**
   - `DB_TYPE`: `local` または `postgres`（省略可、デフォルトは`local`）
   - `DATABASE_URL`: ローカルPostgreSQLの接続文字列
     - 例: `postgresql://postgres:postgres@localhost:5432/mokuhyo_db_test`
   - `SUPABASE_JWKS_URL`: SupabaseのJWKS URL（認証用）
   - `PORT`: 必要に応じて変更（デフォルトは8080）

   **デプロイ時:**
   - `DB_TYPE`: `supabase`
   - `DATABASE_URL`: Supabaseの接続文字列（パスワードを実際の値に置き換え）
   - `SUPABASE_JWKS_URL`: `https://[YOUR-PROJECT-REF].supabase.co/auth/v1/keys`
   - `PORT`: 必要に応じて変更（デフォルトは8080）

### ローカルPostgreSQLのセットアップ

ローカル開発でPostgreSQLを使用する場合:

1. **PostgreSQLをインストール**
   ```bash
   # macOS (Homebrew)
   brew install postgresql@15
   brew services start postgresql@15

   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   ```

2. **データベースを作成**
   ```bash
   # PostgreSQLに接続
   psql postgres

   # データベースとユーザーを作成
   CREATE DATABASE mokuhyo_db_test;
   CREATE USER postgres WITH PASSWORD 'postgres';
   GRANT ALL PRIVILEGES ON DATABASE mokuhyo_db_test TO postgres;
   ```

3. **マイグレーションを実行**
   ```bash
   # SupabaseのマイグレーションファイルをローカルDBに適用
   psql -U postgres -d mokuhyo_db_test -f supabase/migrations/001_initial_schema.sql
   ```

### 注意事項

- `.env` ファイルは **gitにコミットしないでください**（`.gitignore`に含まれています）
- パスワードや機密情報が含まれるため、共有しないでください

## 3. フロントエンド（Next.js）の環境変数

### ファイル場所
```
apps/web/.env.local
```

### 設定内容

```env
# SupabaseプロジェクトのURL
NEXT_PUBLIC_SUPABASE_URL=https://your-project-ref.supabase.co

# Supabase Anon Key（公開キー）
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# バックエンドAPIのベースURL
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### 設定手順

1. `apps/web/.env.local.example` をコピーして `.env.local` を作成
   ```bash
   cd apps/web
   cp .env.local.example .env.local
   ```

2. `.env.local` ファイルを編集して実際の値を設定
   - `NEXT_PUBLIC_SUPABASE_URL`: SupabaseのProject URL
   - `NEXT_PUBLIC_SUPABASE_ANON_KEY`: SupabaseのAnon Key
   - `NEXT_PUBLIC_API_BASE_URL`: 開発環境では `http://localhost:8080`、本番環境ではデプロイしたAPIのURL

### 注意事項

- `.env.local` ファイルは **gitにコミットしないでください**（`.gitignore`に含まれています）
- `NEXT_PUBLIC_` プレフィックスが付いた変数は、ブラウザに公開されます（機密情報は含めないでください）
- Anon Keyは公開キーなので、フロントエンドで使用しても問題ありません（RLSで保護されています）

## 4. 環境変数の確認

### バックエンドの確認
```bash
cd apps/api
go run cmd/server/main.go
# エラーなく起動すればOK
```

### フロントエンドの確認
```bash
cd apps/web
npm run dev
# エラーなく起動すればOK
```

## 5. 本番環境での設定

### デプロイプラットフォーム別の設定方法

#### Fly.io（バックエンド）
```bash
fly secrets set DB_TYPE=supabase
fly secrets set DATABASE_URL=postgresql://...
fly secrets set SUPABASE_JWKS_URL=https://...
fly secrets set PORT=8080
```

#### Vercel / Cloudflare Pages（フロントエンド）
- ダッシュボードの「Environment Variables」セクションで設定
- `NEXT_PUBLIC_SUPABASE_URL`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`
- `NEXT_PUBLIC_API_BASE_URL`

## 6. トラブルシューティング

### データベース接続エラー
- `DB_TYPE` が正しく設定されているか確認（`local`、`postgres`、または`supabase`）
- `DATABASE_URL` のパスワードが正しいか確認
- ローカル開発時: PostgreSQLが起動しているか確認（`brew services list` または `sudo systemctl status postgresql`）
- Supabase使用時: Supabaseのプロジェクトがアクティブか確認
- ファイアウォール設定を確認（SupabaseのIP制限がある場合）

### JWT検証エラー
- `SUPABASE_JWKS_URL` が正しいか確認
- Supabaseの認証設定が有効か確認

### フロントエンドの認証エラー
- `NEXT_PUBLIC_SUPABASE_URL` と `NEXT_PUBLIC_SUPABASE_ANON_KEY` が正しいか確認
- Google認証プロバイダーが有効になっているか確認（Supabase Dashboard > Authentication > Providers）
