# Google認証の設定ガイド

## エラーの原因

エラー `Unsupported provider: provider is not enabled` は、SupabaseでGoogle認証プロバイダーが有効になっていないことが原因です。

## 解決方法

### ステップ1: SupabaseでGoogle認証を有効化

1. **Supabase Dashboardにアクセス**
   - https://supabase.com/dashboard
   - プロジェクトを選択

2. **Authentication設定を開く**
   - 左メニューから「Authentication」をクリック
   - 「Providers」タブを選択

3. **Google認証を有効化**
   - 「Google」プロバイダーを探す
   - 「Enable Google provider」のトグルをONにする

### ステップ2: Google Cloud ConsoleでOAuth認証情報を作成

1. **Google Cloud Consoleにアクセス**
   - https://console.cloud.google.com/
   - Googleアカウントでログイン

2. **プロジェクトを作成（または既存のプロジェクトを選択）**
   - 画面上部のプロジェクト選択ドロップダウンをクリック
   - 「新しいプロジェクト」をクリック
   - プロジェクト名を入力（例: "mokuhyo-driven-test"）
   - 「作成」をクリック

3. **OAuth同意画面を設定**
   - 左メニューから「APIとサービス」→「OAuth同意画面」を選択
   - 「外部」を選択して「作成」をクリック
   - アプリ情報を入力：
     - **アプリ名**: 目標ドリブン・フローティングノート（または任意の名前）
     - **ユーザーサポートメール**: あなたのメールアドレス
     - **デベロッパーの連絡先情報**: あなたのメールアドレス
   - 「保存して次へ」をクリック
   - スコープはデフォルトのままで「保存して次へ」
   - テストユーザーは後で追加できるので「保存して次へ」
   - 「ダッシュボードに戻る」をクリック

4. **認証情報を作成**
   - 左メニューから「APIとサービス」→「認証情報」を選択
   - 画面上部の「+ 認証情報を作成」→「OAuth 2.0 クライアント ID」をクリック
   - **アプリケーションの種類**: 「ウェブアプリケーション」を選択
   - **名前**: 任意の名前（例: "Supabase Auth"）
   - **承認済みの JavaScript 生成元**: 以下を追加
     ```
     https://fyvuqllumshzkrfyytmx.supabase.co
     ```
   - **承認済みのリダイレクト URI**: 以下を追加
     ```
     https://fyvuqllumshzkrfyytmx.supabase.co/auth/v1/callback
     ```
   - 「作成」をクリック

5. **クライアントIDとクライアントシークレットをコピー**
   - 作成されたOAuth 2.0 クライアント IDのダイアログが表示されます
   - **クライアントID**をコピー
   - **クライアントシークレット**をコピー（「シークレットを表示」をクリックして表示）

### ステップ3: Supabaseに認証情報を設定

1. **Supabase Dashboardに戻る**
   - Authentication > Providers > Google

2. **認証情報を入力**
   - **Client ID (Web client)**: Google Cloud ConsoleでコピーしたクライアントIDを貼り付け
   - **Client Secret (Web client)**: Google Cloud Consoleでコピーしたクライアントシークレットを貼り付け

3. **保存**
   - 「Save」ボタンをクリック

### ステップ4: 動作確認

1. **フロントエンドサーバーを再起動**（必要に応じて）
   ```bash
   cd apps/web
   npm run dev
   ```

2. **ブラウザでアクセス**
   - http://localhost:3000 にアクセス
   - 「始める」ボタンをクリック
   - 「Googleでログイン」ボタンをクリック
   - Googleアカウントの選択画面が表示されれば成功

## トラブルシューティング

### エラー: "redirect_uri_mismatch"

**原因**: Google Cloud ConsoleのリダイレクトURIが正しく設定されていない

**解決方法**:
1. Google Cloud Console > 認証情報 > OAuth 2.0 クライアント IDを開く
2. 「承認済みのリダイレクト URI」に以下が含まれているか確認：
   ```
   https://fyvuqllumshzkrfyytmx.supabase.co/auth/v1/callback
   ```
3. 含まれていない場合は追加して保存

### エラー: "access_denied"

**原因**: OAuth同意画面の設定が完了していない

**解決方法**:
1. Google Cloud Console > APIとサービス > OAuth同意画面
2. 設定が完了しているか確認
3. テストユーザーを追加する必要がある場合：
   - 「テストユーザー」タブを開く
   - 「+ ユーザーを追加」をクリック
   - 自分のメールアドレスを追加

### エラー: "invalid_client"

**原因**: クライアントIDまたはクライアントシークレットが間違っている

**解決方法**:
1. Supabase Dashboard > Authentication > Providers > Google
2. クライアントIDとクライアントシークレットが正しく入力されているか確認
3. Google Cloud Consoleで再度コピーして貼り付け

### ローカル開発環境での注意

**開発環境（localhost）でテストする場合**:
- Google Cloud Consoleの「承認済みのリダイレクト URI」に以下も追加：
  ```
  http://localhost:3000/auth/callback
  ```
- ただし、Supabase経由で認証する場合は、SupabaseのリダイレクトURIのみで動作します

## まとめ

1. ✅ SupabaseでGoogle認証プロバイダーを有効化
2. ✅ Google Cloud ConsoleでOAuth認証情報を作成
3. ✅ SupabaseにクライアントIDとシークレットを設定
4. ✅ リダイレクトURIが正しく設定されているか確認

これでGoogle認証が動作するはずです！
