# 404エラーの診断と解決方法

## 現在の状況

1. **Next.jsの開発サーバーは起動している**（2つのプロセスが実行中）
2. **しかし、404エラーが返ってくる**
3. **「missing required error components, refreshing...」というメッセージが表示される**

## 原因

Next.jsのビルドが不完全で、必要なエラーコンポーネントが欠けています。また、複数のNext.jsプロセスが起動しているため、競合している可能性があります。

## 解決手順

### ステップ1: すべてのNext.jsプロセスを停止

```bash
# すべてのNext.jsプロセスを確認
ps aux | grep "next dev" | grep -v grep

# プロセスIDを確認して、すべて停止
kill -9 93281 47000
# （実際のプロセスIDに置き換えてください）
```

または、ターミナルで`Ctrl+C`を押して開発サーバーを停止してください。

### ステップ2: ビルドキャッシュを完全にクリア

```bash
cd apps/web
rm -rf .next
rm -rf node_modules/.cache
```

### ステップ3: エラーコンポーネントを追加（必要に応じて）

Next.js 14では、エラーコンポーネントが必要な場合があります。以下のファイルを作成してください：

`apps/web/app/error.tsx`:
```tsx
'use client'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-4">エラーが発生しました</h2>
        <button
          onClick={() => reset()}
          className="bg-blue-600 text-white px-4 py-2 rounded"
        >
          再試行
        </button>
      </div>
    </div>
  )
}
```

`apps/web/app/not-found.tsx`:
```tsx
export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-4">ページが見つかりません</h2>
        <a href="/" className="text-blue-600 underline">
          ホームに戻る
        </a>
      </div>
    </div>
  )
}
```

### ステップ4: 開発サーバーを再起動

```bash
cd apps/web
npm run dev
```

### ステップ5: ブラウザのキャッシュをクリア

1. ブラウザの開発者ツール（F12）を開く
2. 「Application」タブを開く
3. 「Storage」セクションで「Clear site data」をクリック
4. すべてにチェックを入れて「Clear site data」をクリック

## 確認すべき情報

問題を解決するために、以下の情報を確認してください：

1. **ターミナルのログ**: Next.jsの開発サーバーを起動した際に、どのようなメッセージが表示されていますか？
2. **ブラウザのコンソール**: ブラウザの開発者ツール（F12）のコンソールタブに、どのようなエラーが表示されていますか？
3. **環境変数**: `apps/web/.env.local`ファイルは存在しますか？内容は何ですか？
4. **アクセスしているURL**: ブラウザでどのURLにアクセスしていますか？（`http://localhost:3000`、`http://localhost:3000/login`など）

## それでも解決しない場合

完全にクリーンアップして再インストール：

```bash
cd apps/web
rm -rf .next
rm -rf node_modules
rm package-lock.json
npm install
npm run dev
```
