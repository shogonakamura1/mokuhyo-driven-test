'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useAuth } from '@/hooks/useAuth'
import { api } from '@/lib/api'
import type { Project } from '@/types'

interface SidebarProps {
  project: Project
}

export default function Sidebar({ project }: SidebarProps) {
  const { signOut } = useAuth()
  const [projects, setProjects] = useState<Project[]>([])
  const [loadingProjects, setLoadingProjects] = useState(true)

  useEffect(() => {
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
  }, [])

  const handleLogout = async () => {
    await signOut()
  }

  return (
    <div className="w-64 bg-white border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <h2 className="text-lg font-semibold text-gray-900">{project.title}</h2>
        {project.description && (
          <p className="text-sm text-gray-600 mt-1">{project.description}</p>
        )}
      </div>
      <div className="flex-1 p-4 space-y-4 overflow-auto">
        <div>
          <h3 className="text-sm font-semibold text-gray-700 mb-2">過去のツリー</h3>
          {loadingProjects ? (
            <div className="text-xs text-gray-500">読み込み中...</div>
          ) : projects.length === 0 ? (
            <div className="text-xs text-gray-500">まだ作成したツリーがありません</div>
          ) : (
            <div className="space-y-1">
              {projects.map(item => {
                const isCurrent = item.id === project.id
                return (
                  <Link
                    key={item.id}
                    href={`/p/${item.id}`}
                    className={`block rounded-md border px-2 py-1 text-xs transition-colors ${
                      isCurrent
                        ? 'border-blue-400 bg-blue-50 text-blue-700'
                        : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50'
                    }`}
                  >
                    <div className="font-medium truncate">{item.title}</div>
                    <div className="text-[10px] text-gray-500">
                      最終更新: {new Date(item.updated_at).toLocaleString('ja-JP')}
                    </div>
                  </Link>
                )
              })}
            </div>
          )}
        </div>
      </div>
      <div className="p-4 border-t border-gray-200">
        <button
          onClick={handleLogout}
          className="w-full text-gray-600 hover:text-gray-900 transition-colors"
        >
          ログアウト
        </button>
      </div>
    </div>
  )
}
