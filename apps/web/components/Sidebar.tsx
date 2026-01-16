'use client'

import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import type { Project } from '@/types'

interface SidebarProps {
  project: Project
  onSave: () => void
}

export default function Sidebar({ project, onSave }: SidebarProps) {
  const router = useRouter()
  const { signOut } = useAuth()

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
      <div className="flex-1 p-4 space-y-2">
        <button
          onClick={onSave}
          className="w-full bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
        >
          保存
        </button>
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
