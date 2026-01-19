# redirect_uri_mismatch エラーの最終解決方法

## 問題の原因

`@react-oauth/google`の`useGoogleLogin`が、デフォルトで現在のページのURLをリダイレクトURIとして使用しますが、これがGoogle Cloud Consoleに登録されているURIと一致していない可能性があります。

## 解決方法

### 方法1: ポップアップモードを使用（推奨・実装済み）

`useGoogleLogin`に`ux_mode: 'popup'`を追加しました。これにより、リダイレクトURIの問題を回避できます。

**変更内容**:
```typescript
const signInWithGoogle = useGoogleLogin({
  flow: 'auth-code',
  ux_mode: 'popup', // ポップアップモードを使用
  onSuccess: async (codeResponse) => {
    // ...
  },
})
```

**次のステップ**:
1. 開発サーバーを再起動（必要に応じて）
2. ブラウザで`http://localhost:3001/login`にアクセス
3. 「Googleでログイン」をクリック
4. ポップアップウィンドウでGoogle認証が表示されるはずです

### 方法2: リダイレクトモードを使用する場合

もしリダイレクトモードを使用したい場合は、`redirect_uri`を明示的に指定する必要があります：

```typescript
const signInWithGoogle = useGoogleLogin({
  flow: 'auth-code',
  ux_mode: 'redirect',
  redirect_uri: window.location.origin, // 現在のオリジンを明示的に指定
  onSuccess: async (codeResponse) => {
    // ...
  },
})
```

この場合、Google Cloud Consoleに以下のURIを追加する必要があります：
- `http://localhost:3000`
- `http://localhost:3001`

## 確認事項

### 1. ブラウザで正しいポートにアクセスしているか

- ✅ `http://localhost:3001`（Next.jsが起動しているポート）
- ❌ `http://localhost:3000`（別のプロセスが使用中）

### 2. Google Cloud Consoleの設定

以下のURIが登録されていることを確認：
- `http://localhost:3000`
- `http://localhost:3001`

### 3. ポップアップがブロックされていないか

ブラウザのポップアップブロッカーが有効になっている場合、ポップアップが表示されない可能性があります。

**解決方法**:
- ブラウザの設定でポップアップを許可する
- または、リダイレクトモードを使用する（方法2）

## トラブルシューティング

### ポップアップが表示されない場合

1. **ブラウザのポップアップブロッカーを確認**
   - アドレスバーの右側にポップアップブロックのアイコンが表示されていないか確認
   - ポップアップを許可する

2. **コンソールでエラーを確認**
   - ブラウザの開発者ツール（F12）を開く
   - コンソールタブでエラーメッセージを確認

3. **リダイレクトモードに切り替える**
   - 方法2を試す

### エラーが続く場合

1. **Google Cloud Consoleの設定を再確認**
   - 設定が正しく保存されているか確認
   - 数分待ってから再度試す（設定の反映に時間がかかる場合があります）

2. **ブラウザのキャッシュをクリア**
   - ブラウザの開発者ツール（F12）を開く
   - 「Application」タブ → 「Storage」 → 「Clear site data」

3. **プライベートモードで試す**
   - ブラウザの拡張機能が干渉している可能性があります
