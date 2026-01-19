# Google OAuth redirect_uri_mismatch エラーの完全解決方法

## 現在の状況

- Next.jsはポート3001で起動している
- ブラウザは`http://localhost:3000`にアクセスしている
- Google Cloud Consoleには`http://localhost:3000`と`http://localhost:3001`の両方が登録されている
- しかし、まだ`redirect_uri_mismatch`エラーが発生している

## 原因

`@react-oauth/google`の`useGoogleLogin`は、ポップアップモードでも**承認済みのJavaScript生成元**の設定が必要です。

## 解決方法

### ステップ1: Google Cloud Consoleで「承認済みのJavaScript生成元」を設定

1. **Google Cloud Consoleにアクセス**
   - https://console.cloud.google.com/
   - プロジェクトを選択

2. **認証情報を開く**
   - 左メニューから「APIとサービス」→「認証情報」を選択
   - OAuth 2.0 クライアント IDをクリック

3. **「承認済みのJavaScript生成元」に追加**
   - 「承認済みのJavaScript生成元」セクションを開く
   - 「+ URI を追加」をクリック
   - 以下のURIを追加：
     ```
     http://localhost:3000
     http://localhost:3001
     ```
   - 「保存」をクリック

4. **「承認済みのリダイレクトURI」も確認**
   - 以下のURIが登録されていることを確認：
     ```
     http://localhost:3000
     http://localhost:3001
     ```

### ステップ2: ポート3000を使用しているプロセスを停止

```bash
# ポート3000を使用しているプロセスを確認
lsof -i :3000

# プロセスを停止（PIDを確認してから）
kill -9 <PID>
```

### ステップ3: Next.jsをポート3000で起動

```bash
cd apps/web
npm run dev
```

これで、Next.jsはポート3000で起動するはずです。

### ステップ4: ブラウザでアクセス

1. `http://localhost:3000/login`にアクセス
2. 「Googleでログイン」をクリック
3. ポップアップウィンドウでGoogle認証が表示されるはずです

## 重要な設定ポイント

### Google Cloud Consoleで設定すべきURI

**承認済みのJavaScript生成元**（必須）:
- `http://localhost:3000`
- `http://localhost:3001`

**承認済みのリダイレクトURI**（念のため）:
- `http://localhost:3000`
- `http://localhost:3001`

## トラブルシューティング

### エラーが続く場合

1. **Google Cloud Consoleの設定を再確認**
   - 「承認済みのJavaScript生成元」に`http://localhost:3000`が追加されているか確認
   - 設定を保存してから、数分待つ（設定の反映に時間がかかる場合があります）

2. **ブラウザのキャッシュをクリア**
   - ブラウザの開発者ツール（F12）を開く
   - 「Application」タブ → 「Storage」 → 「Clear site data」

3. **プライベートモードで試す**
   - ブラウザの拡張機能が干渉している可能性があります

4. **コンソールでエラーを確認**
   - ブラウザの開発者ツール（F12）を開く
   - 「Console」タブでエラーメッセージを確認

## 確認チェックリスト

- [ ] Google Cloud Consoleの「承認済みのJavaScript生成元」に`http://localhost:3000`が追加されている
- [ ] Google Cloud Consoleの「承認済みのリダイレクトURI」に`http://localhost:3000`が追加されている
- [ ] ポート3000を使用しているプロセスを停止した
- [ ] Next.jsをポート3000で起動した
- [ ] ブラウザで`http://localhost:3000/login`にアクセスしている
- [ ] ブラウザのキャッシュをクリアした
