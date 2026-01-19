# ポート3000で起動する方法

## 問題

ブラウザは`http://localhost:3000`にアクセスしていますが、Next.jsはポート3001で起動しています。これが`redirect_uri_mismatch`エラーの原因です。

## 解決方法

### ステップ1: ポート3000を使用しているプロセスを停止

```bash
# ポート3000を使用しているプロセスを確認
lsof -i :3000

# プロセスを停止（例：PIDが47001の場合）
kill -9 47001
```

### ステップ2: Next.jsの開発サーバーを再起動

```bash
cd apps/web
npm run dev
```

これで、Next.jsはポート3000で起動するはずです。

### ステップ3: ブラウザでアクセス

`http://localhost:3000/login`にアクセスして、「Googleでログイン」をクリックしてください。

## 確認事項

### Google Cloud Consoleの設定

以下のURIが登録されていることを確認してください：

**承認済みのJavaScript生成元**:
- `http://localhost:3000`

**承認済みのリダイレクトURI**:
- `http://localhost:3000`
- `http://localhost:3001`（念のため）

## トラブルシューティング

### ポート3000がまだ使用されている場合

```bash
# すべてのNodeプロセスを確認
ps aux | grep node

# ポート3000を使用しているプロセスを強制終了
lsof -ti :3000 | xargs kill -9
```

### それでもポート3000で起動しない場合

環境変数でポートを明示的に指定：

```bash
cd apps/web
PORT=3000 npm run dev
```
