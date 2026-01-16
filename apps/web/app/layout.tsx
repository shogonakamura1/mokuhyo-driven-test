import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: '目標ドリブン・フローティングノート',
  description: '目標に向かって思考を深掘りするアプリ',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  )
}
