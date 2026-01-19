# トラブルシューティングガイド

## ログイン画面が「読み込み中」のまま表示される

### 原因1: 環境変数が設定されていない

**確認方法**:
1. ブラウザの開発者ツール（F12）を開く
2. コンソールタブを確認
3. エラーメッセージを確認

**解決方法**:
1. `apps/web/.env.local`ファイルを作成（存在しない場合）
2. 以下の環境変数を設定：

```env
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

3. Next.jsの開発サーバーを再起動：
```bash
cd apps/web
npm run dev
```

### 原因2: Google Client IDが正しく設定されていない

**確認方法**:
1. ブラウザの開発者ツール（F12）を開く
2. コンソールで以下を実行：
```javascript
console.log(process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID)
```

**解決方法**:
1. `.env.local`ファイルの`NEXT_PUBLIC_GOOGLE_CLIENT_ID`が正しいか確認
2. Google Cloud ConsoleでクライアントIDを確認
3. 開発サーバーを再起動

### 原因3: GoogleOAuthProviderが正しく初期化されていない

**確認方法**:
1. ブラウザの開発者ツール（F12）を開く
2. コンソールでエラーを確認
3. ネットワークタブでGoogle OAuthのリクエストを確認

**解決方法**:
1. `layout.tsx`で環境変数が正しく読み込まれているか確認
2. 開発サーバーを再起動

### 原因4: localStorageへのアクセスエラー

**確認方法**:
1. ブラウザの開発者ツール（F12）を開く
2. コンソールでエラーを確認

**解決方法**:
1. ブラウザのプライベートモードで試す
2. ブラウザのキャッシュをクリア
3. ローカルストレージをクリア：
```javascript
localStorage.clear()
```

## その他のよくある問題

### エラー: "Failed to exchange code for token"

**原因**: バックエンドの認証エンドポイントが正しく動作していない

**解決方法**:
1. バックエンドサーバーが起動しているか確認
2. `apps/api/.env`ファイルに以下が設定されているか確認：
```env
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000
```
3. バックエンドのログを確認

### エラー: "redirect_uri_mismatch"

**原因**: Google Cloud ConsoleのリダイレクトURIが正しく設定されていない

**解決方法**:
1. Google Cloud Console > 認証情報 > OAuth 2.0 クライアント IDを開く
2. 「承認済みのリダイレクト URI」に以下を追加：
```
http://localhost:3000
```

### エラー: "user not found"

**原因**: データベースにユーザーが作成されていない

**解決方法**:
1. データベースの`users`テーブルが作成されているか確認
2. マイグレーションを実行：
```bash
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/001_initial_schema.sql
psql -U postgres -d mokuhyo_db_test -f supabase/migrations/002_add_users_table.sql
```

## デバッグのヒント

1. **ブラウザのコンソールを確認**: エラーメッセージが表示されることが多い
2. **ネットワークタブを確認**: APIリクエストが正しく送信されているか確認
3. **バックエンドのログを確認**: サーバー側でエラーが発生していないか確認
4. **環境変数を確認**: `.env.local`ファイルが正しく設定されているか確認
5. **開発サーバーを再起動**: 環境変数の変更後は再起動が必要
