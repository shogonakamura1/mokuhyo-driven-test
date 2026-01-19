# Next.js 404エラーの解決方法

## エラーの内容

```
（インデックス）:1  Failed to load resource: the server responded with a status of 404 (Not Found)
（インデックス）:5  GET http://localhost:3000/ 404 (Not Found)
check @ （インデックス）:5
setTimeout
check @ （インデックス）:10
```

## 原因

Next.jsの開発サーバーが起動していますが、正しくルートパス（`/`）を処理できていません。`check`関数は、おそらくブラウザの拡張機能やService Workerが原因です。

## 解決方法

### 方法1: Next.jsの開発サーバーを完全に再起動（推奨）

```bash
cd apps/web

# 1. 開発サーバーを停止（Ctrl+C）

# 2. ビルドキャッシュを削除
rm -rf .next

# 3. node_modulesのキャッシュも削除（必要に応じて）
rm -rf node_modules/.cache

# 4. 開発サーバーを再起動
npm run dev
```

### 方法2: ポートを変更して起動

既存のプロセスが干渉している可能性があります。

```bash
cd apps/web

# ポート3001で起動
PORT=3001 npm run dev
```

その後、ブラウザで`http://localhost:3001`にアクセスしてください。

### 方法3: ブラウザのキャッシュとService Workerをクリア

1. **ブラウザの開発者ツール（F12）を開く**
2. **「Application」タブを開く**
3. **「Storage」セクションで「Clear site data」をクリック**
4. **すべてにチェックを入れて「Clear site data」をクリック**

または、ブラウザのコンソールで以下を実行：

```javascript
// Service Workerを登録解除
navigator.serviceWorker.getRegistrations().then(function(registrations) {
  for(let registration of registrations) {
    registration.unregister();
  }
});

// キャッシュをクリア
caches.keys().then(function(names) {
  for (let name of names) caches.delete(name);
});

// ローカルストレージをクリア
localStorage.clear();
```

### 方法4: プライベートモードで試す

ブラウザの拡張機能が干渉している可能性があります。

1. プライベートモード（シークレットモード）を開く
2. `http://localhost:3000`にアクセス

### 方法5: 別のブラウザで試す

ブラウザ固有の問題の可能性があります。

## 確認事項

### 1. Next.jsの開発サーバーが正しく起動しているか

ターミナルで以下のようなメッセージが表示されているか確認：

```
  ▲ Next.js 14.2.0
  - Local:        http://localhost:3000
  - ready started server on 0.0.0.0:3000, url: http://localhost:3000
```

### 2. ビルドエラーがないか確認

ターミナルにエラーメッセージが表示されていないか確認してください。

### 3. 環境変数が正しく設定されているか

`apps/web/.env.local`ファイルが存在し、正しく設定されているか確認：

```env
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## トラブルシューティング

### エラーが続く場合

1. **完全にクリーンアップ**:
```bash
cd apps/web
rm -rf .next
rm -rf node_modules/.cache
npm run dev
```

2. **Node.jsのバージョンを確認**:
```bash
node --version
# Next.js 14にはNode.js 18以上が必要
```

3. **依存関係を再インストール**:
```bash
cd apps/web
rm -rf node_modules
rm package-lock.json
npm install
npm run dev
```

### `check`関数が繰り返し実行される場合

これは、ブラウザの拡張機能やService Workerが原因の可能性が高いです。

**解決方法**:
1. ブラウザの拡張機能を一時的に無効にする
2. プライベートモードで試す
3. Service Workerを登録解除（上記の手順参照）

## 推奨される手順

1. **開発サーバーを停止**（Ctrl+C）
2. **`.next`ディレクトリを削除**: `rm -rf .next`
3. **開発サーバーを再起動**: `npm run dev`
4. **ブラウザのキャッシュとService Workerをクリア**
5. **プライベートモードで試す**
6. **それでも解決しない場合は、別のポートで起動**: `PORT=3001 npm run dev`
