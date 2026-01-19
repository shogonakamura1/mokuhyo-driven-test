# データベースマイグレーションガイド

## 1. データベースの作成

まず、PostgreSQLに接続してデータベースを作成します。

```bash
# PostgreSQLに接続（デフォルトのpostgresデータベースに接続）
psql -U postgres

# データベースを作成
CREATE DATABASE mokuhyo_db_test;

# ユーザーを作成（既に存在する場合はスキップ）
CREATE USER postgres WITH PASSWORD 'postgres';

# 権限を付与
GRANT ALL PRIVILEGES ON DATABASE mokuhyo_db_test TO postgres;

# 接続を終了
\q
```

## 2. マイグレーションの実行

### 方法1: psqlコマンドで直接実行（推奨）

```bash
# 1つ目のマイグレーション（初期スキーマ）を実行
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/001_initial_schema.sql

# 2つ目のマイグレーション（usersテーブル）を実行
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/002_add_users_table.sql
```

### 方法2: psqlに接続して実行

```bash
# PostgreSQLに接続
psql -U postgres -d mokuhyo_db_test

# マイグレーションファイルの内容をコピー&ペーストして実行
# または、\iコマンドで実行
\i supabase/migrations/001_initial_schema.sql
\i supabase/migrations/002_add_users_table.sql

# 接続を終了
\q
```

## 3. マイグレーションの確認

マイグレーションが正しく実行されたか確認します。

```bash
# PostgreSQLに接続
psql -U postgres -d mokuhyo_db_test

# テーブル一覧を確認
\dt

# usersテーブルの構造を確認
\d users

# projectsテーブルの構造を確認
\d projects

# 外部キー制約を確認
SELECT
    tc.table_name, 
    kcu.column_name, 
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM 
    information_schema.table_constraints AS tc 
    JOIN information_schema.key_column_usage AS kcu
      ON tc.constraint_name = kcu.constraint_name
    JOIN information_schema.constraint_column_usage AS ccu
      ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY';

# 接続を終了
\q
```

## 4. トラブルシューティング

### エラー: "database does not exist"

データベースが作成されていません。上記の手順1でデータベースを作成してください。

### エラー: "relation already exists"

テーブルが既に存在しています。マイグレーションファイルは`CREATE TABLE IF NOT EXISTS`を使用しているので、このエラーは通常発生しませんが、外部キー制約の追加でエラーが発生する可能性があります。

外部キー制約が既に存在する場合：

```sql
-- 既存の制約を確認
SELECT constraint_name 
FROM information_schema.table_constraints 
WHERE table_name = 'projects' 
  AND constraint_type = 'FOREIGN KEY';

-- 既存の制約を削除（必要に応じて）
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_user_id_fkey;
ALTER TABLE user_settings DROP CONSTRAINT IF EXISTS user_settings_user_id_fkey;

-- 再度マイグレーションを実行
```

### エラー: "permission denied"

ユーザーに適切な権限がありません。

```sql
-- データベースに接続
psql -U postgres -d mokuhyo_db_test

-- 権限を付与
GRANT ALL PRIVILEGES ON DATABASE mokuhyo_db_test TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;

-- デフォルト権限を設定
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO postgres;
```

## 5. マイグレーションの順序

マイグレーションは以下の順序で実行する必要があります：

1. **001_initial_schema.sql**: 基本的なテーブル（projects, nodes, edges, user_settings）を作成
2. **002_add_users_table.sql**: usersテーブルを作成し、外部キー制約を追加

## 6. データベースのリセット（開発時）

開発中にデータベースをリセットしたい場合：

```bash
# PostgreSQLに接続
psql -U postgres

# データベースを削除（注意：すべてのデータが削除されます）
DROP DATABASE IF EXISTS mokuhyo_db_test;

# データベースを再作成
CREATE DATABASE mokuhyo_db_test;
GRANT ALL PRIVILEGES ON DATABASE mokuhyo_db_test TO postgres;

# 接続を終了
\q

# マイグレーションを再実行
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/001_initial_schema.sql
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/002_add_users_table.sql
```

## 7. 環境変数の確認

マイグレーション実行後、環境変数が正しく設定されているか確認してください。

`apps/api/.env`:
```env
DB_TYPE=local
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mokuhyo_db_test
```

## 8. 動作確認

マイグレーション実行後、APIサーバーを起動して動作確認：

```bash
cd apps/api
go run cmd/server/main.go
```

正常に起動すれば、マイグレーションは成功しています。
