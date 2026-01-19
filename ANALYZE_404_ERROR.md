# localhost:3000/ 404エラーの分析と解決方法

## エラーの内容

```
localhost:3000/ 404 (Not Found)
（インデックス）:5  GET http://localhost:3000/ 404 (Not Found)
check @ （インデックス）:5
setTimeout
check @ （インデックス）:10
```

このエラーは、`check`という関数が`setTimeout`を使って繰り返し`localhost:3000/`にリクエストを送信していることを示しています。

## 考えられる原因

### 1. Next.jsの開発サーバーが正しく起動していない

**確認方法**:
```bash
cd apps/web
npm run dev
```

ターミナルに以下のようなメッセージが表示されるか確認：
```
  ▲ Next.js 14.2.0
  - Local:        http://localhost:3000
  - ready started server on 0.0.0.0:3000
```

**解決方法**:
- 開発サーバーが起動していない場合は起動
- 既に起動している場合は、一度停止（Ctrl+C）して再起動

### 2. ブラウザの拡張機能が干渉している

**確認方法**:
1. ブラウザのプライベートモード（シークレットモード）で試す
2. 拡張機能を無効にして試す

**解決方法**:
- プライベートモードで動作する場合は、拡張機能が原因
- 拡張機能を一時的に無効にするか、別のブラウザで試す

### 3. Service Workerが残っている

**確認方法**:
1. ブラウザの開発者ツール（F12）を開く
2. 「Application」タブを開く
3. 「Service Workers」を確認

**解決方法**:
```javascript
// ブラウザのコンソールで実行
navigator.serviceWorker.getRegistrations().then(function(registrations) {
  for(let registration of registrations) {
    registration.unregister();
  }
});
```

### 4. ブラウザのキャッシュが古い

**解決方法**:
1. ブラウザの開発者ツール（F12）を開く
2. ネットワークタブを開く
3. 「Disable cache」にチェックを入れる
4. ページをハードリロード（Ctrl+Shift+R または Cmd+Shift+R）

### 5. Next.jsのビルドキャッシュが壊れている

**解決方法**:
```bash
cd apps/web
rm -rf .next
npm run dev
```

## 推奨される解決手順

### ステップ1: 開発サーバーを再起動

```bash
cd apps/web
# 開発サーバーを停止（Ctrl+C）
rm -rf .next
npm run dev
```

### ステップ2: ブラウザのキャッシュとService Workerをクリア

1. ブラウザの開発者ツール（F12）を開く
2. 「Application」タブを開く
3. 「Storage」セクションで「Clear site data」をクリック
4. すべてにチェックを入れて「Clear site data」をクリック

### ステップ3: プライベートモードで試す

プライベートモード（シークレットモード）で`http://localhost:3000`にアクセスして、動作するか確認

### ステップ4: 別のブラウザで試す

別のブラウザ（Chrome、Firefox、Safariなど）で試して、ブラウザ固有の問題か確認

## デバッグのヒント

### ブラウザのコンソールで詳細を確認

1. ブラウザの開発者ツール（F12）を開く
2. 「Console」タブを開く
3. エラーメッセージのスタックトレースを確認
4. 「Sources」タブで`check`関数がどこで定義されているか確認

### ネットワークタブでリクエストを確認

1. ブラウザの開発者ツール（F12）を開く
2. 「Network」タブを開く
3. `localhost:3000/`へのリクエストを確認
4. リクエストの「Initiator」列で、どのスクリプトがリクエストを送信しているか確認

## よくある原因と解決方法

### 原因: ブラウザ拡張機能（例: React DevTools、Redux DevTools）

**解決方法**: 拡張機能を一時的に無効にする

### 原因: Service Workerが古い状態のまま

**解決方法**: Service Workerを登録解除（上記の手順参照）

### 原因: Next.jsの開発サーバーが正しく起動していない

**解決方法**: 
1. ポート3000が使用されていないか確認：`lsof -i :3000`
2. 使用されている場合は終了：`kill -9 <PID>`
3. 開発サーバーを再起動

### 原因: ビルドエラーが発生している

**解決方法**:
```bash
cd apps/web
npm run build
```
ビルドエラーが表示された場合は、そのエラーを修正

## それでも解決しない場合

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
