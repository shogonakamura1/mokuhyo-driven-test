import Link from 'next/link'

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-4 text-gray-900">ページが見つかりません</h2>
        <Link href="/" className="text-blue-600 underline hover:text-blue-800">
          ホームに戻る
        </Link>
      </div>
    </div>
  )
}
