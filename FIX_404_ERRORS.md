# Next.js 404エラーの解決方法

## エラーの内容

以下のエラーが発生しています：
- `layout.css:1 Failed to load resource: the server responded with a status of 404 (Not Found)`
- `main-app.js:1 Failed to load resource: the server responded with a status of 404 (Not Found)`
- `app-pages-internals.js:1 Failed to load resource: the server responded with a status of 404 (Not Found)`

## 原因

これらのエラーは、Next.jsの開発サーバーが正しくビルドファイルを生成していない、または古いビルドキャッシュが残っている場合に発生します。

## 解決方法

### 方法1: ビルドキャッシュをクリアして再起動（推奨）

```bash
cd apps/web

# .nextディレクトリを削除
rm -rf .next

# node_modulesもクリアする場合（必要に応じて）
# rm -rf node_modules
# npm install

# 開発サーバーを再起動
npm run dev
```

### 方法2: ポートを変更して再起動

既存のプロセスがポート3000を使用している可能性があります。

```bash
cd apps/web

# ポート3000を使用しているプロセスを確認
lsof -i :3000

# プロセスを終了（必要に応じて）
# kill -9 <PID>

# 別のポートで起動
PORT=3001 npm run dev
```

### 方法3: 完全にクリーンアップして再ビルド

```bash
cd apps/web

# すべてのキャッシュとビルドファイルを削除
rm -rf .next
rm -rf node_modules/.cache

# 開発サーバーを再起動
npm run dev
```

## 確認事項

### 1. 開発サーバーが正しく起動しているか確認

ターミナルで以下のようなメッセージが表示されているか確認：

```
  ▲ Next.js 14.2.0
  - Local:        http://localhost:3000
  - ready started server on 0.0.0.0:3000, url: http://localhost:3000
```

### 2. ブラウザのキャッシュをクリア

1. ブラウザの開発者ツール（F12）を開く
2. ネットワークタブを開く
3. 「Disable cache」にチェックを入れる
4. ページをリロード（Ctrl+Shift+R または Cmd+Shift+R）

### 3. 環境変数が正しく設定されているか確認

`apps/web/.env.local`ファイルが存在し、正しく設定されているか確認：

```env
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## トラブルシューティング

### エラーが続く場合

1. **別のブラウザで試す**: ブラウザ固有の問題の可能性
2. **プライベートモードで試す**: ブラウザの拡張機能が干渉している可能性
3. **ターミナルのログを確認**: 開発サーバーの起動時にエラーが表示されていないか確認

### ビルドエラーが発生する場合

```bash
cd apps/web

# ビルドを試行してエラーを確認
npm run build
```

ビルドエラーが表示された場合は、そのエラーメッセージに従って修正してください。

## よくある原因

1. **開発サーバーが正常に起動していない**: ポートが既に使用されている、または起動に失敗している
2. **古いビルドキャッシュ**: `.next`ディレクトリが古い状態のまま
3. **ファイルの変更が反映されていない**: ホットリロードが機能していない
4. **環境変数の問題**: `.env.local`ファイルが正しく読み込まれていない

## 推奨される手順

1. 開発サーバーを停止（Ctrl+C）
2. `.next`ディレクトリを削除
3. 開発サーバーを再起動
4. ブラウザのキャッシュをクリア
5. ページをリロード

これで通常は解決します。
