# データベース接続エラーの解決方法

## 問題の原因

エラー `no route to host` は、以下の原因が考えられます：

1. **IPv6接続の問題**: Supabaseの直接接続はIPv6を使用するため、IPv4しか使えない環境では接続できません
2. **DATABASE_URLの形式が間違っている**: 接続文字列の形式が正しくない可能性があります

## 解決方法

### 方法1: SupabaseのPooler接続を使用（推奨）

SupabaseのPooler接続（セッションモード）を使用すると、IPv4経由で接続できます。

1. **Supabase Dashboardにアクセス**
   - https://supabase.com/dashboard
   - プロジェクトを選択

2. **Pooler接続文字列を取得**
   - Settings > Database > Connection string
   - **Connection pooling** タブを選択
   - **Session mode** を選択
   - **URI** をコピー

3. **`.env`ファイルを更新**
   ```env
   DATABASE_URL=postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres
   ```

   例：
   ```env
   DATABASE_URL=postgresql://postgres.fyvuqllumshzkrfyytmx:your-password@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres
   ```

### 方法2: DATABASE_URLの形式を確認・修正

現在の`.env`ファイルの`DATABASE_URL`が以下の形式になっているか確認してください：

**正しい形式:**
```env
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres
```

**間違った形式の例:**
```env
# ❌ 間違い: postgresql:が抜けている
DATABASE_URL=postgres:***@07Sho516@db.fyvuqllumshzkrfyytmx.supabase.co:5432/postgres

# ❌ 間違い: パスワードの位置が間違っている
DATABASE_URL=postgresql:***@password@host:port/database
```

**正しい形式の例:**
```env
# ✅ 正しい
DATABASE_URL=postgresql://postgres:your-actual-password@db.fyvuqllumshzkrfyytmx.supabase.co:5432/postgres
```

### 方法3: コードでIPv4接続を強制（既に実装済み）

コードを修正してIPv4接続を強制するようにしました。これで直接接続でもIPv4経由で接続できるようになります。

## 確認手順

1. **`.env`ファイルの確認**
   ```bash
   cd apps/api
   cat .env | grep DATABASE_URL
   ```

2. **接続文字列の形式を確認**
   - `postgresql://` で始まっているか
   - `postgres:[PASSWORD]@` の形式になっているか
   - パスワードに特殊文字が含まれている場合は、URLエンコードが必要な場合があります

3. **Supabase Dashboardで確認**
   - Settings > Database > Connection string
   - **Direct connection** と **Connection pooling** の両方を確認
   - Pooler接続の方が安定している場合があります

## 推奨される設定

**`.env`ファイル（Pooler接続を使用）:**
```env
# Pooler接続（セッションモード）- IPv4対応
DATABASE_URL=postgresql://postgres.fyvuqllumshzkrfyytmx:[PASSWORD]@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres

# または直接接続（IPv4強制コードで対応）
DATABASE_URL=postgresql://postgres:[PASSWORD]@db.fyvuqllumshzkrfyytmx.supabase.co:5432/postgres

SUPABASE_JWKS_URL=https://fyvuqllumshzkrfyytmx.supabase.co/auth/v1/keys
PORT=8080
```

## トラブルシューティング

### まだ接続できない場合

1. **パスワードに特殊文字が含まれている場合**
   - URLエンコードが必要です
   - 例: `@` → `%40`, `#` → `%23`

2. **SupabaseのIP制限**
   - Settings > Database > Network Restrictions
   - 接続元のIPが許可されているか確認

3. **SupabaseのIP Ban**
   - Settings > Database > Network Bans
   - IPがbanされていないか確認

4. **接続テスト**
   ```bash
   # psqlで接続テスト（インストールされている場合）
   psql "postgresql://postgres:[PASSWORD]@db.fyvuqllumshzkrfyytmx.supabase.co:5432/postgres"
   ```
