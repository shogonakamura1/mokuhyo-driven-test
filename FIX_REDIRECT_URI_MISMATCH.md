# redirect_uri_mismatch エラーの解決方法

## エラーの内容

**Error 400: redirect_uri_mismatch**

このエラーは、Google OAuthの設定で、リダイレクトURIが一致していないことを示しています。

## 原因

Next.jsの開発サーバーがポート3001で起動していますが、Google Cloud Consoleの設定ではポート3000のリダイレクトURIが登録されている可能性があります。

## 解決方法

### 方法1: Google Cloud Consoleでポート3001のリダイレクトURIを追加（推奨）

1. **Google Cloud Consoleにアクセス**
   - https://console.cloud.google.com/
   - プロジェクトを選択

2. **認証情報を開く**
   - 左メニューから「APIとサービス」→「認証情報」を選択
   - OAuth 2.0 クライアント IDをクリック

3. **リダイレクトURIを追加**
   - 「承認済みのリダイレクト URI」セクションを開く
   - 以下のURIを追加：
     ```
     http://localhost:3001
     http://localhost:3001/auth/callback
     ```
   - 「保存」をクリック

4. **ブラウザで再度ログインを試す**
   - `http://localhost:3001/login`にアクセス
   - 「Googleでログイン」をクリック

### 方法2: ポート3000を使用しているプロセスを停止して、ポート3000で起動

```bash
# ポート3000を使用しているプロセスを停止
kill -9 47001

# Next.jsの開発サーバーを再起動（ポート3000で起動するはず）
cd apps/web
npm run dev
```

その後、ブラウザで`http://localhost:3000`にアクセスしてください。

### 方法3: 環境変数でリダイレクトURIを指定

`apps/api/.env`ファイルに以下を追加：

```env
GOOGLE_REDIRECT_URL=http://localhost:3001
```

または、ポート3000で起動する場合：

```env
GOOGLE_REDIRECT_URL=http://localhost:3000
```

## 確認すべきリダイレクトURI

Google Cloud Consoleで以下のURIが登録されているか確認してください：

**開発環境用**:
- `http://localhost:3000`
- `http://localhost:3001`（現在のポート）

**本番環境用**（デプロイ時）:
- `https://your-production-domain.com`
- `https://your-production-domain.com/auth/callback`

## 注意事項

- リダイレクトURIは**完全一致**する必要があります
- 末尾のスラッシュ（`/`）の有無も重要です
- `http://`と`https://`は別物として扱われます

## トラブルシューティング

### エラーが続く場合

1. **Google Cloud ConsoleでリダイレクトURIを確認**
   - 設定が保存されているか確認
   - 数分待ってから再度試す（設定の反映に時間がかかる場合があります）

2. **ブラウザのキャッシュをクリア**
   - ブラウザの開発者ツール（F12）を開く
   - 「Application」タブ → 「Storage」 → 「Clear site data」

3. **プライベートモードで試す**
   - ブラウザの拡張機能が干渉している可能性があります
