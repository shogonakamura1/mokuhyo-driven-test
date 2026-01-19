# クイックスタートガイド

## サーバー起動までの手順

### 前提条件
- ✅ Supabaseプロジェクトの作成が完了している
- ✅ 環境変数ファイル（`.env` / `.env.local`）の設定が完了している

---

## ステップ1: データベースマイグレーションの実行

Supabaseでデータベーススキーマを作成します。

1. **Supabase Dashboardにアクセス**
   - https://supabase.com/dashboard にログイン
   - プロジェクトを選択

2. **SQL Editorを開く**
   - 左メニューから「SQL Editor」をクリック

3. **マイグレーションファイルを実行**
   - 「New query」をクリック
   - 以下のファイルの内容をコピー＆ペースト:
     ```
     supabase/migrations/001_initial_schema.sql
     ```
   - 「Run」ボタンをクリックして実行
   - 成功メッセージが表示されればOK

4. **確認**
   - 左メニューから「Table Editor」を開く
   - 以下のテーブルが作成されていることを確認:
     - `projects`
     - `nodes`
     - `edges`
     - `user_settings`
     - `snapshots`

---

## ステップ2: Google認証の設定（Supabase）

1. **Supabase Dashboardで認証設定**
   - 左メニューから「Authentication」→「Providers」を開く
   - 「Google」を選択
   - 「Enable Google provider」をONにする

2. **Google OAuth認証情報の取得**
   - [Google Cloud Console](https://console.cloud.google.com/)にアクセス
   - プロジェクトを作成（または既存のプロジェクトを選択）
   - 「APIとサービス」→「認証情報」を開く
   - 「認証情報を作成」→「OAuth 2.0 クライアント ID」を選択
   - アプリケーションの種類: 「ウェブアプリケーション」
   - 承認済みのリダイレクト URI に以下を追加:
     ```
     https://[YOUR-PROJECT-REF].supabase.co/auth/v1/callback
     ```
   - クライアントIDとクライアントシークレットをコピー

3. **Supabaseに認証情報を設定**
   - Supabase Dashboard > Authentication > Providers > Google
   - Client ID（Web client）とClient Secret（Web client）を入力
   - 「Save」をクリック

---

## ステップ3: バックエンドサーバーの起動

### 3.1 依存関係のインストール（初回のみ）

```bash
cd apps/api
go mod download
```

### 3.2 サーバーの起動

```bash
# apps/apiディレクトリにいることを確認
cd apps/api

# サーバーを起動
go run cmd/server/main.go
```

**成功時の表示:**
```
Server starting on port 8080
```

**エラーが出た場合:**
- `.env`ファイルが正しく設定されているか確認
- `DATABASE_URL`のパスワードが正しいか確認
- Supabaseのプロジェクトがアクティブか確認

### 3.3 動作確認

別のターミナルで以下を実行して、サーバーが起動しているか確認:

```bash
curl http://localhost:8080/health
```

**期待されるレスポンス:**
```json
{"status":"ok"}
```

---

## ステップ4: フロントエンドサーバーの起動

### 4.1 依存関係のインストール（初回のみ）

```bash
cd apps/web
npm install
```

### 4.2 開発サーバーの起動

```bash
# apps/webディレクトリにいることを確認
cd apps/web

# 開発サーバーを起動
npm run dev
```

**成功時の表示:**
```
- ready started server on 0.0.0.0:3000, url: http://localhost:3000
- event compiled client and server successfully
```

### 4.3 ブラウザでアクセス

ブラウザで以下にアクセス:
```
http://localhost:3000
```

---

## ステップ5: 動作確認

### 5.1 ログイン
1. ブラウザで `http://localhost:3000` にアクセス
2. 「始める」ボタンをクリック
3. Googleアカウントでログイン

### 5.2 プロジェクト作成
1. ログイン後、目標のタイトルを入力（3〜20文字）
2. 「作成」ボタンをクリック

### 5.3 ツリー編集
1. メイン編集画面が表示される
2. ノードをダブルクリックして編集
3. **Enterキー**: 編集確定 → 兄弟ノード追加
4. **Tabキー**: 編集確定 → 子ノード追加
5. **Shift+Enter**: 編集確定のみ
6. **Escキー**: 編集キャンセル
7. **右クリック**: ノード削除

---

## トラブルシューティング

### バックエンドが起動しない

**エラー: `DATABASE_URL is required`**
- `.env`ファイルが`apps/api`ディレクトリに存在するか確認
- 環境変数が正しく設定されているか確認

**エラー: `failed to connect to database`**
- `DATABASE_URL`のパスワードが正しいか確認
- Supabaseのプロジェクトがアクティブか確認
- データベース接続文字列の形式が正しいか確認

**エラー: `SUPABASE_JWKS_URL is required`**
- `.env`ファイルに`SUPABASE_JWKS_URL`が設定されているか確認

### フロントエンドが起動しない

**エラー: `Missing Supabase environment variables`**
- `.env.local`ファイルが`apps/web`ディレクトリに存在するか確認
- `NEXT_PUBLIC_SUPABASE_URL`と`NEXT_PUBLIC_SUPABASE_ANON_KEY`が設定されているか確認

**エラー: `Module not found`**
- `npm install`を実行して依存関係をインストール

### ログインできない

- Google認証プロバイダーがSupabaseで有効になっているか確認
- Google Cloud ConsoleのリダイレクトURIが正しく設定されているか確認
- ブラウザのコンソールでエラーメッセージを確認

### APIリクエストが失敗する

- バックエンドサーバーが起動しているか確認（`http://localhost:8080/health`）
- ブラウザの開発者ツール（Networkタブ）でエラーを確認
- CORSエラーが出ている場合は、バックエンドのCORS設定を確認

---

## 次のステップ

サーバーが正常に起動したら、以下を試してみてください:

1. **プロジェクトの作成**: 目標を入力してプロジェクトを作成
2. **ノードの追加**: EnterキーやTabキーでノードを追加
3. **ノードの編集**: ダブルクリックで編集
4. **ノードの削除**: 右クリックで削除
5. **保存**: サイドバーの「保存」ボタンで保存

---

## 開発時の便利なコマンド

### バックエンド
```bash
# サーバー起動
cd apps/api && go run cmd/server/main.go

# ビルド
cd apps/api && go build ./cmd/server

# テスト
cd apps/api && go test ./...
```

### フロントエンド
```bash
# 開発サーバー起動
cd apps/web && npm run dev

# ビルド
cd apps/web && npm run build

# 本番サーバー起動
cd apps/web && npm start

# リンター
cd apps/web && npm run lint
```

---

## サーバーを停止する

- **バックエンド**: ターミナルで `Ctrl + C`
- **フロントエンド**: ターミナルで `Ctrl + C`
