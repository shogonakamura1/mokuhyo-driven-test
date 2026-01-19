import type { Metadata } from 'next'
import './globals.css'
import { GoogleOAuthProvider } from '@react-oauth/google'

export const metadata: Metadata = {
  title: '目標ドリブン・フローティングノート',
  description: '目標に向かって思考を深掘りするアプリ',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const googleClientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || ''

  return (
    <html lang="ja">
      <body>
        {googleClientId ? (
          <GoogleOAuthProvider clientId={googleClientId}>
            {children}
          </GoogleOAuthProvider>
        ) : (
          <div className="min-h-screen flex items-center justify-center">
            <div className="text-red-600">
              <p className="font-bold">エラー: Google Client IDが設定されていません</p>
              <p className="text-sm mt-2">.env.localファイルにNEXT_PUBLIC_GOOGLE_CLIENT_IDを設定してください</p>
            </div>
          </div>
        )}
      </body>
    </html>
  )
}
