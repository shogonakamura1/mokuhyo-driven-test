# redirect_uri_mismatch エラーの重要な修正ポイント

## 問題の原因

`@react-oauth/google`の`useGoogleLogin`は、**ポップアップモードでも「承認済みのJavaScript生成元」の設定が必要**です。

## 必須の設定

### Google Cloud Consoleで設定すべき項目

1. **「承認済みのJavaScript生成元」**（重要！）
   - `http://localhost:3000`
   - `http://localhost:3001`

2. **「承認済みのリダイレクトURI」**（念のため）
   - `http://localhost:3000`
   - `http://localhost:3001`

## 手順

### 1. Google Cloud Consoleで設定を追加

1. Google Cloud Consoleにアクセス
2. 「APIとサービス」→「認証情報」を選択
3. OAuth 2.0 クライアント IDをクリック
4. **「承認済みのJavaScript生成元」セクションを開く**
5. 「+ URI を追加」をクリック
6. 以下を追加：
   ```
   http://localhost:3000
   ```
7. もう一度「+ URI を追加」をクリック
8. 以下を追加：
   ```
   http://localhost:3001
   ```
9. 「保存」をクリック

### 2. ポート3000を使用しているプロセスを停止

```bash
# ポート3000を使用しているプロセスを停止
kill -9 9826

# 確認
lsof -i :3000
```

### 3. Next.jsを再起動

```bash
cd apps/web
npm run dev
```

これで、Next.jsはポート3000で起動するはずです。

### 4. ブラウザでアクセス

`http://localhost:3000/login`にアクセスして、「Googleでログイン」をクリックしてください。

## 重要なポイント

- **「承認済みのJavaScript生成元」は必須です**（ポップアップモードでも必要）
- 設定を保存してから、数分待つ（設定の反映に時間がかかる場合があります）
- ブラウザのキャッシュをクリアする

## 確認方法

設定が正しく反映されているか確認：

1. Google Cloud Consoleで「承認済みのJavaScript生成元」を確認
2. `http://localhost:3000`と`http://localhost:3001`の両方が表示されているか確認
3. 設定を保存してから、数分待つ
4. ブラウザのキャッシュをクリア
5. 再度ログインを試す
