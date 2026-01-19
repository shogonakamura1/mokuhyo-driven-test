# データベース環境変数設定ガイド

このプロジェクトでは、クリーンアーキテクチャに基づき、`DB_TYPE`環境変数でデータベース実装を切り替えることができます。

## 環境変数の設定

### ファイル場所
```
apps/api/.env
```

## 設定例

### 1. ローカル開発時（ローカルPostgreSQLを使用）

```env
# データベースタイプ
# local または postgres: ローカルPostgreSQLを使用
# デフォルト: local（省略可能）
DB_TYPE=local

# データベース接続文字列（ローカルPostgreSQL）
# 形式: postgresql://user:password@localhost:5432/dbname
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mokuhyo_db_test

# Google OAuth設定
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000

# APIサーバーのポート番号（デフォルト: 8080）
PORT=8080
```

### 2. デプロイ時（Supabaseを使用）

```env
# データベースタイプ
DB_TYPE=supabase

# データベース接続文字列（Supabase PostgreSQL）
# 形式: postgresql://postgres:password@db.project-ref.supabase.co:5432/postgres
DATABASE_URL=postgresql://postgres:your-password@db.your-project-ref.supabase.co:5432/postgres

# Google OAuth設定
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=https://your-production-domain.com

# APIサーバーのポート番号（デフォルト: 8080）
PORT=8080
```

## 環境変数の説明

### DB_TYPE
- **値**: `local`、`postgres`、または`supabase`
- **デフォルト**: `local`（省略可能）
- **説明**: 使用するデータベース実装を指定します
  - `local` または `postgres`: ローカルPostgreSQLを使用（開発時）
  - `supabase`: Supabaseを使用（デプロイ時）

### DATABASE_URL
- **説明**: PostgreSQLデータベースへの接続文字列
- **ローカル開発時**: `postgresql://user:password@localhost:5432/dbname`
- **Supabase使用時**: `postgresql://postgres:password@db.project-ref.supabase.co:5432/postgres`

### SUPABASE_JWKS_URL
- **説明**: SupabaseのJWKS（JSON Web Key Set）URL。JWTトークンの検証に使用されます
- **形式**: `https://your-project-ref.supabase.co/auth/v1/keys`
- **注意**: ローカル開発時でも必要です（認証はSupabaseを使用）

### PORT
- **説明**: APIサーバーがリッスンするポート番号
- **デフォルト**: `8080`
- **省略可能**: はい

## ローカルPostgreSQLのセットアップ

### 1. PostgreSQLのインストール

#### macOS (Homebrew)
```bash
brew install postgresql@15
brew services start postgresql@15
```

#### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### 2. データベースの作成

```bash
# PostgreSQLに接続
psql postgres

# データベースとユーザーを作成
CREATE DATABASE mokuhyo_db_test;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE mokuhyo_db_test TO postgres;

# 接続を終了
\q
```

### 3. マイグレーションの実行

```bash
# SupabaseのマイグレーションファイルをローカルDBに適用
psql -U postgres -d mokuhyo_db_test_test -f supabase/migrations/001_initial_schema.sql
```

または、PostgreSQLに接続して直接実行:

```bash
psql -U postgres -d mokuhyo_db_test
```

```sql
-- supabase/migrations/001_initial_schema.sql の内容をコピー&ペースト
```

## 動作確認

### 1. 環境変数の確認

```bash
cd apps/api
cat .env
```

### 2. サーバーの起動

```bash
cd apps/api
go run cmd/server/main.go
```

起動時に以下のようなログが表示されます:

```
Using local PostgreSQL database implementation
Server starting on port 8080
```

または（Supabase使用時）:

```
Using Supabase database implementation
Server starting on port 8080
```

### 3. ヘルスチェック

```bash
curl http://localhost:8080/health
```

正常な場合、以下のレスポンスが返ります:

```json
{"status":"ok"}
```

## トラブルシューティング

### データベース接続エラー

1. **`DB_TYPE`が正しく設定されているか確認**
   ```bash
   echo $DB_TYPE  # または .envファイルを確認
   ```

2. **ローカル開発時: PostgreSQLが起動しているか確認**
   ```bash
   # macOS
   brew services list
   
   # Linux
   sudo systemctl status postgresql
   ```

3. **接続文字列が正しいか確認**
   - ユーザー名、パスワード、データベース名が正しいか
   - ポート番号が5432か（デフォルト）

4. **データベースが存在するか確認**
   ```bash
   psql -U postgres -l
   ```

### JWT検証エラー

- `SUPABASE_JWKS_URL`が正しいか確認
- Supabaseのプロジェクトがアクティブか確認
- ネットワーク接続を確認

### ポートが既に使用されている

```bash
# ポート8080を使用しているプロセスを確認
lsof -i :8080

# 別のポートを使用する場合、.envファイルでPORTを変更
PORT=8081
```

## デプロイ時の設定

### Fly.io

```bash
fly secrets set DB_TYPE=supabase
fly secrets set DATABASE_URL=postgresql://...
fly secrets set SUPABASE_JWKS_URL=https://...
fly secrets set PORT=8080
```

### Docker

```dockerfile
ENV DB_TYPE=supabase
ENV DATABASE_URL=postgresql://...
ENV SUPABASE_JWKS_URL=https://...
ENV PORT=8080
```

### Kubernetes

```yaml
env:
  - name: DB_TYPE
    value: "supabase"
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: url
  - name: SUPABASE_JWKS_URL
    value: "https://..."
  - name: PORT
    value: "8080"
```
