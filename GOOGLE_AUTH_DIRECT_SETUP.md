# Google認証直接使用の設定ガイド

このプロジェクトでは、Supabase経由ではなく、Google OAuthを直接使用して認証を行います。

## 1. Google Cloud ConsoleでOAuth認証情報を作成

### ステップ1: Google Cloud Consoleにアクセス

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. Googleアカウントでログイン

### ステップ2: プロジェクトを作成（または既存のプロジェクトを選択）

1. 画面上部のプロジェクト選択ドロップダウンをクリック
2. 「新しいプロジェクト」をクリック
3. プロジェクト名を入力（例: "mokuhyo-driven-test"）
4. 「作成」をクリック

### ステップ3: OAuth同意画面を設定

1. 左メニューから「APIとサービス」→「OAuth同意画面」を選択
2. 「外部」を選択して「作成」をクリック
3. アプリ情報を入力：
   - **アプリ名**: 目標ドリブン・フローティングノート（または任意の名前）
   - **ユーザーサポートメール**: あなたのメールアドレス
   - **デベロッパーの連絡先情報**: あなたのメールアドレス
4. 「保存して次へ」をクリック
5. スコープはデフォルトのままで「保存して次へ」
6. テストユーザーを追加（必要に応じて）：
   - 「テストユーザー」タブを開く
   - 「+ ユーザーを追加」をクリック
   - 自分のメールアドレスを追加
7. 「保存して次へ」をクリック
8. 「ダッシュボードに戻る」をクリック

### ステップ4: OAuth 2.0 クライアントIDを作成

1. 左メニューから「APIとサービス」→「認証情報」を選択
2. 画面上部の「+ 認証情報を作成」→「OAuth 2.0 クライアント ID」をクリック
3. **アプリケーションの種類**: 「ウェブアプリケーション」を選択
4. **名前**: 任意の名前（例: "Mokuhyo Web Client"）
5. **承認済みの JavaScript 生成元**: 以下を追加
   ```
   http://localhost:3000
   https://your-production-domain.com
   ```
6. **承認済みのリダイレクト URI**: 以下を追加
   ```
   http://localhost:3000
   https://your-production-domain.com
   ```
7. 「作成」をクリック
8. **クライアントID**と**クライアントシークレット**をコピーして保存

## 2. 環境変数の設定

### バックエンド（Go API）

`apps/api/.env`ファイルに以下を設定：

```env
# Google OAuth設定
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000

# データベース設定
DB_TYPE=local
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/mokuhyo_db_test

# APIサーバーのポート
PORT=8080
```

### フロントエンド（Next.js）

`apps/web/.env.local`ファイルに以下を設定：

```env
# Google OAuth設定
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com

# バックエンドAPIのベースURL
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## 3. データベースのセットアップ

### usersテーブルの作成

マイグレーションファイルを実行：

```bash
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/002_add_users_table.sql
```

または、PostgreSQLに接続して直接実行：

```sql
-- usersテーブルを作成
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  google_id TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL,
  name TEXT NOT NULL,
  picture TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS users_google_id_idx ON users(google_id);
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);

-- projectsテーブルとの外部キー制約を追加
ALTER TABLE projects 
  ADD CONSTRAINT projects_user_id_fkey 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- user_settingsテーブルとの外部キー制約を追加
ALTER TABLE user_settings
  ADD CONSTRAINT user_settings_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
```

## 4. 動作確認

### バックエンドの起動

```bash
cd apps/api
go run cmd/server/main.go
```

### フロントエンドの起動

```bash
cd apps/web
npm run dev
```

### 認証フローの確認

1. ブラウザで `http://localhost:3000` にアクセス
2. 「Googleでログイン」ボタンをクリック
3. Googleアカウントの選択画面が表示されれば成功
4. ログイン後、`/input` ページにリダイレクトされる

## 5. トラブルシューティング

### エラー: "redirect_uri_mismatch"

**原因**: Google Cloud ConsoleのリダイレクトURIが正しく設定されていない

**解決方法**:
1. Google Cloud Console > 認証情報 > OAuth 2.0 クライアント IDを開く
2. 「承認済みのリダイレクト URI」に以下が含まれているか確認：
   ```
   http://localhost:3000
   ```
3. 含まれていない場合は追加して保存

### エラー: "access_denied"

**原因**: OAuth同意画面の設定が完了していない、またはテストユーザーが追加されていない

**解決方法**:
1. Google Cloud Console > APIとサービス > OAuth同意画面
2. 設定が完了しているか確認
3. テストユーザーを追加（開発中の場合）

### エラー: "invalid_client"

**原因**: クライアントIDまたはクライアントシークレットが間違っている

**解決方法**:
1. `.env`ファイルの`GOOGLE_CLIENT_ID`と`GOOGLE_CLIENT_SECRET`が正しいか確認
2. Google Cloud Consoleで再度コピーして貼り付け

### エラー: "user not found"

**原因**: データベースにユーザーが作成されていない

**解決方法**:
1. データベースの`users`テーブルが作成されているか確認
2. マイグレーションファイルを実行
3. バックエンドのログを確認してエラーがないか確認

## 6. 本番環境での設定

### Google Cloud Console

1. OAuth同意画面を「本番」に変更
2. 承認済みのリダイレクトURIに本番URLを追加：
   ```
   https://your-production-domain.com
   ```

### 環境変数

**バックエンド**:
```env
GOOGLE_CLIENT_ID=your-production-client-id
GOOGLE_CLIENT_SECRET=your-production-client-secret
GOOGLE_REDIRECT_URL=https://your-production-domain.com
```

**フロントエンド**:
```env
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-production-client-id
NEXT_PUBLIC_API_BASE_URL=https://your-api-domain.com
```

## まとめ

1. ✅ Google Cloud ConsoleでOAuth認証情報を作成
2. ✅ 環境変数を設定
3. ✅ データベースにusersテーブルを作成
4. ✅ バックエンドとフロントエンドを起動
5. ✅ 認証フローを確認

これでGoogle認証が直接使用できるようになりました！
