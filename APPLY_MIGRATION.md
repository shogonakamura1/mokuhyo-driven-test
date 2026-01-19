# データベースマイグレーションの適用方法

## 問題

データベースのチェック制約`nodes_content_check`が`char_length(content) between 1 and 200`となっているため、空文字列を許可していません。

## 解決方法

新しいマイグレーションファイル`003_allow_empty_node_content.sql`を作成しました。このマイグレーションをデータベースに適用してください。

## 手順

### 1. PostgreSQLに接続

```bash
psql -d mokuhyo_db_test -U your_username
```

または、環境変数`DATABASE_URL`が設定されている場合：

```bash
psql $DATABASE_URL
```

### 2. マイグレーションを実行

```sql
-- 既存の制約を削除
ALTER TABLE nodes
DROP CONSTRAINT IF EXISTS nodes_content_check;

-- 新しい制約を追加（空文字列を許可、最大200文字）
ALTER TABLE nodes
ADD CONSTRAINT nodes_content_check CHECK (char_length(content) <= 200);
```

### 3. 確認

```sql
-- 制約が正しく設定されているか確認
SELECT conname, pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conrelid = 'nodes'::regclass
AND conname = 'nodes_content_check';
```

出力例：
```
nodes_content_check | CHECK ((char_length(content) <= 200))
```

## または、マイグレーションファイルを直接実行

```bash
psql -d mokuhyo_db_test -U your_username -f supabase/migrations/003_allow_empty_node_content.sql
```

または：

```bash
psql $DATABASE_URL -f supabase/migrations/003_allow_empty_node_content.sql
```

## 確認方法

マイグレーション適用後、空文字列でノードを作成できることを確認してください。
