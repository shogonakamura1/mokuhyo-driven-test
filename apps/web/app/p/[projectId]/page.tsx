'use client'

import { useEffect, useState } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { api } from '@/lib/api'
import Sidebar from '@/components/Sidebar'
import TreeCanvas from '@/components/TreeCanvas'
import type { TreeResponse, Node, Edge } from '@/types'

export default function ProjectPage() {
  const { user, loading } = useAuth()
  const router = useRouter()
  const params = useParams()
  const projectId = params.projectId as string

  const [tree, setTree] = useState<TreeResponse | null>(null)
  const [loadingTree, setLoadingTree] = useState(true)
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)
  const [editingNodeId, setEditingNodeId] = useState<string | null>(null)

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login')
      return
    }

    if (user && projectId) {
      loadTree()
    }
  }, [user, loading, projectId, router])

  const loadTree = async () => {
    try {
      setLoadingTree(true)
      const data = await api.get<TreeResponse>(`/v1/projects/${projectId}/tree`)
      setTree(data)
    } catch (error) {
      console.error('Failed to load tree:', error)
      alert('ツリーの読み込みに失敗しました')
    } finally {
      setLoadingTree(false)
    }
  }

  const handleNodeCreate = async (
    content: string,
    parentNodeId: string | null,
    relation: string = 'neutral',
    orderIndex?: number,
    question?: string
  ) => {
    try {
      const response = await api.post<{ node: Node; edge: Edge }>(
        `/v1/projects/${projectId}/nodes`,
        {
          content,
          parent_node_id: parentNodeId || null,
          relation,
          order_index: orderIndex,
          question,
        }
      )
      await loadTree()
      // ツリーが再読み込みされた後、新しいノードを選択して編集モードにする
      setSelectedNodeId(response.node.id)
      setEditingNodeId(response.node.id)
      return response.node
    } catch (error) {
      console.error('Failed to create node:', error)
      throw error
    }
  }

  const handleNodeUpdate = async (nodeId: string, content: string) => {
    try {
      await api.patch(`/v1/projects/${projectId}/nodes/${nodeId}`, { content })
      await loadTree()
      // 更新後、編集モードを解除
      if (editingNodeId === nodeId) {
        setEditingNodeId(null)
      }
    } catch (error) {
      console.error('Failed to update node:', error)
      throw error
    }
  }

  const handleNodeDelete = async (nodeId: string) => {
    if (!confirm('このノードとその子ノードを削除しますか？')) {
      return
    }
    try {
      await api.delete(`/v1/projects/${projectId}/nodes/${nodeId}`)
      await loadTree()
      if (selectedNodeId === nodeId) {
        setSelectedNodeId(null)
      }
      if (editingNodeId === nodeId) {
        setEditingNodeId(null)
      }
    } catch (error) {
      console.error('Failed to delete node:', error)
      alert('削除に失敗しました')
    }
  }

  const handleEdgeUpdate = async (edgeId: string, relation: string, relationLabel?: string) => {
    try {
      await api.patch(`/v1/projects/${projectId}/edges/${edgeId}`, {
        relation,
        relation_label: relationLabel,
      })
      await loadTree()
    } catch (error) {
      console.error('Failed to update edge:', error)
      throw error
    }
  }

  if (loading || loadingTree) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-600">読み込み中...</div>
      </div>
    )
  }

  if (!tree) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-red-600">ツリーの読み込みに失敗しました</div>
      </div>
    )
  }

  return (
    <div className="flex h-screen bg-gray-50">
      <Sidebar
        project={tree.project}
      />
      <div className="flex-1 overflow-auto">
        <TreeCanvas
          nodes={tree.nodes ?? []}
          edges={tree.edges ?? []}
          selectedNodeId={selectedNodeId}
          editingNodeId={editingNodeId}
          onSelectNode={setSelectedNodeId}
          onSetEditingNodeId={setEditingNodeId}
          onNodeCreate={handleNodeCreate}
          onNodeUpdate={handleNodeUpdate}
          onNodeDelete={handleNodeDelete}
          onEdgeUpdate={handleEdgeUpdate}
        />
      </div>
    </div>
  )
}
