# Supabase APIキーの使用箇所

## APIキーの種類と用途

SupabaseのAPI Keys画面には2種類のキーが表示されています：

### 1. Publishable Key（公開キー / anon key）

**用途**: フロントエンド（Next.js）で使用

**使用箇所**:
- **ファイル**: `apps/web/.env.local`
- **環境変数名**: `NEXT_PUBLIC_SUPABASE_ANON_KEY`
- **使用されるコード**: `apps/web/lib/supabase.ts`

**具体的な使用例**:

```typescript
// apps/web/lib/supabase.ts
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
export const supabase = createClient(supabaseUrl, supabaseAnonKey)
```

**このキーでできること**:
- Google認証（OAuth）
- Supabase Authのセッション管理
- Row Level Security (RLS) が有効なテーブルへのアクセス

**セキュリティ**:
- ✅ ブラウザに公開されても安全（RLSで保護されているため）
- ✅ フロントエンドで使用可能
- ⚠️ ただし、RLSポリシーが正しく設定されていることが前提

### 2. Secret Key（秘密キー / service_role key）

**このプロジェクトでの使用**: **使用していません**

**理由**:
- このプロジェクトでは、Go APIがJWT検証を直接行うため、SupabaseのSDKを使用していません
- バックエンドではSupabaseのAPIキーではなく、JWTトークンの検証のみを行います

**通常の使用例（参考）**:
- サーバーサイドでSupabaseのSDKを使用する場合
- RLSをバイパスしてデータベースに直接アクセスする場合
- 管理機能やバッチ処理など

## 設定方法

### 1. Publishable Keyの設定

1. **Supabase Dashboardでキーをコピー**
   - Settings > API Keys
   - **Publishable key** の「default」のAPI KEYをコピー
   - コピーボタン（📋アイコン）をクリック

2. **`.env.local`ファイルに設定**
   ```bash
   cd apps/web
   ```
   
   `.env.local`ファイルを開いて、以下を設定：
   ```env
   NEXT_PUBLIC_SUPABASE_URL=https://fyvuqllumshzkrfyytmx.supabase.co
   NEXT_PUBLIC_SUPABASE_ANON_KEY=sb_publishable_8n7P-XGu5eTqcuwAbaj-BQ_t7I2T...
   ```

3. **確認**
   - フロントエンドのコードで使用されているか確認：
     ```bash
     grep -r "NEXT_PUBLIC_SUPABASE_ANON_KEY" apps/web/
     ```

### 2. Secret Keyについて

このプロジェクトでは使用しませんが、将来的に必要になった場合：

- **使用しない理由**: Go APIがJWT検証を直接行うため
- **もし使用する場合**: サーバーサイドでのみ使用（環境変数に設定、gitにコミットしない）

## 実際の使用フロー

### 認証フロー

1. **フロントエンド（Next.js）**
   ```typescript
   // apps/web/lib/supabase.ts
   // Publishable Keyを使用してSupabaseクライアントを初期化
   export const supabase = createClient(supabaseUrl, supabaseAnonKey)
   
   // apps/web/hooks/useAuth.ts
   // Google認証を実行
   await supabase.auth.signInWithOAuth({ provider: 'google' })
   ```

2. **認証成功後**
   - SupabaseがJWTトークンを発行
   - フロントエンドがトークンを取得
   - APIリクエスト時にトークンを送信

3. **バックエンド（Go API）**
   ```go
   // apps/api/pkg/auth/jwt.go
   // JWTトークンを検証（Supabase JWKSを使用）
   // APIキーは使用しない
   ```

## トラブルシューティング

### Publishable Keyが正しく設定されていない場合

**エラー**: `Missing Supabase environment variables`

**解決方法**:
1. `.env.local`ファイルが`apps/web`ディレクトリに存在するか確認
2. `NEXT_PUBLIC_SUPABASE_ANON_KEY`が正しく設定されているか確認
3. キーの値が完全にコピーされているか確認（途中で切れていないか）

### 認証が動作しない場合

1. **Publishable Keyが正しいか確認**
   - Supabase Dashboard > Settings > API Keys
   - 最新のキーをコピー

2. **Google認証プロバイダーが有効か確認**
   - Supabase Dashboard > Authentication > Providers
   - Googleが有効になっているか確認

3. **リダイレクトURIが正しく設定されているか確認**
   - Google Cloud Console > OAuth 2.0 クライアント ID
   - 承認済みのリダイレクト URIに以下が含まれているか：
     ```
     https://fyvuqllumshzkrfyytmx.supabase.co/auth/v1/callback
     ```

## まとめ

| キーの種類 | 使用箇所 | 環境変数名 | このプロジェクトでの使用 |
|-----------|---------|-----------|----------------------|
| **Publishable Key** | フロントエンド | `NEXT_PUBLIC_SUPABASE_ANON_KEY` | ✅ 使用中 |
| **Secret Key** | バックエンド（通常） | - | ❌ 使用しない |

**重要なポイント**:
- Publishable Keyはフロントエンドの`.env.local`ファイルで使用
- Secret Keyはこのプロジェクトでは使用しない
- 両方のキーとも、gitにコミットしないように注意
