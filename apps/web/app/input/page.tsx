'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { api } from '@/lib/api'
import type { Project } from '@/types'

export default function InputPage() {
  const { user, loading } = useAuth()
  const router = useRouter()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [projects, setProjects] = useState<Project[]>([])
  const [loadingProjects, setLoadingProjects] = useState(true)

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login')
    }
  }, [user, loading, router])

  useEffect(() => {
    if (loading || !user) return

    const loadProjects = async () => {
      try {
        setLoadingProjects(true)
        const response = await api.get<{ projects: Project[] }>('/v1/projects')
        const sorted = [...response.projects].sort(
          (a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
        )
        setProjects(sorted)
      } catch (error) {
        console.error('Failed to load projects:', error)
        setProjects([])
      } finally {
        setLoadingProjects(false)
      }
    }

    loadProjects()
  }, [user, loading])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim() || title.length < 3 || title.length > 20) {
      alert('タイトルは3文字以上20文字以下で入力してください')
      return
    }

    setSubmitting(true)
    try {
      const response = await api.post<{ project: Project }>('/v1/projects', {
        title: title.trim(),
        description: description.trim() || undefined,
      })
      router.push(`/p/${response.project.id}`)
    } catch (error) {
      console.error('Failed to create project:', error)
      alert('プロジェクトの作成に失敗しました')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-600">読み込み中...</div>
      </div>
    )
  }

  if (!user) {
    return null
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white rounded-lg shadow-md p-8 space-y-6">
        <div>
          <h2 className="text-lg font-semibold mb-3">過去のツリー</h2>
          {loadingProjects ? (
            <div className="text-sm text-gray-500">読み込み中...</div>
          ) : projects.length === 0 ? (
            <div className="text-sm text-gray-500">まだ作成したツリーがありません</div>
          ) : (
            <div className="space-y-2">
              {projects.map(project => (
                <Link
                  key={project.id}
                  href={`/p/${project.id}`}
                  className="block rounded-lg border border-gray-200 px-3 py-2 hover:border-blue-300 hover:bg-blue-50 transition-colors"
                >
                  <div className="text-sm font-medium text-gray-900">{project.title}</div>
                  <div className="text-xs text-gray-500">
                    最終更新: {new Date(project.updated_at).toLocaleString('ja-JP')}
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
        <h1 className="text-2xl font-bold text-center mb-6">目標を入力</h1>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
              タイトル <span className="text-red-500">*</span>
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="例: TOEIC 800点"
              minLength={3}
              maxLength={20}
              required
            />
            <p className="mt-1 text-sm text-gray-500">
              {title.length}/20文字
            </p>
          </div>
          <div>
            <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
              説明（任意）
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              rows={3}
              placeholder="例: 3ヶ月で達成"
            />
          </div>
          <button
            type="submit"
            disabled={submitting || !title.trim()}
            className="w-full bg-blue-600 text-white px-4 py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {submitting ? '作成中...' : '作成'}
          </button>
        </form>
      </div>
    </div>
  )
}
