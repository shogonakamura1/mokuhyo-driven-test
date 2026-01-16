import Link from 'next/link'

export default function HomePage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <h1 className="text-5xl font-bold text-gray-900 mb-6">
          目標ドリブン・フローティングノート
        </h1>
        <p className="text-xl text-gray-700 mb-8">
          目標に向かって思考を深掘りし、具体的な行動に落とし込む
        </p>
        <Link
          href="/login"
          className="inline-block bg-blue-600 text-white px-8 py-3 rounded-lg text-lg font-semibold hover:bg-blue-700 transition-colors"
        >
          始める
        </Link>
      </div>
    </div>
  )
}
