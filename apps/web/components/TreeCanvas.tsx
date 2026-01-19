'use client'

import { useCallback, useMemo } from 'react'
import NodeCard from './NodeCard'
import type { Node, Edge, RelationType } from '@/types'

interface TreeCanvasProps {
  nodes: Node[]
  edges: Edge[]
  selectedNodeId: string | null
  editingNodeId: string | null
  onSelectNode: (nodeId: string | null) => void
  onSetEditingNodeId: (nodeId: string | null) => void
  onNodeCreate: (content: string, parentNodeId: string | null, relation?: string, orderIndex?: number) => Promise<Node>
  onNodeUpdate: (nodeId: string, content: string) => Promise<void>
  onNodeDelete: (nodeId: string) => Promise<void>
  onEdgeUpdate: (edgeId: string, relation: string, relationLabel?: string) => Promise<void>
}

export default function TreeCanvas({
  nodes,
  edges,
  selectedNodeId,
  editingNodeId,
  onSelectNode,
  onSetEditingNodeId,
  onNodeCreate,
  onNodeUpdate,
  onNodeDelete,
  onEdgeUpdate,
}: TreeCanvasProps) {

  // Build tree structure
  const treeStructure = useMemo(() => {
    const nodeMap = new Map(nodes.map(node => [node.id, node]))
    const childrenMap = new Map<string, Array<{ node: Node; edge: Edge }>>()

    edges.forEach(edge => {
      const parentId = edge.parent_node_id || 'root'
      if (!childrenMap.has(parentId)) {
        childrenMap.set(parentId, [])
      }
      const node = nodeMap.get(edge.child_node_id)
      if (node) {
        childrenMap.get(parentId)!.push({ node, edge })
      }
    })

    // Sort children by order_index
    childrenMap.forEach(children => {
      children.sort((a, b) => a.edge.order_index - b.edge.order_index)
    })

    return { nodeMap, childrenMap }
  }, [nodes, edges])

  const handleNodeClick = (nodeId: string) => {
    if (editingNodeId === nodeId) return
    onSelectNode(nodeId)
  }

  const handleNodeDoubleClick = (nodeId: string) => {
    onSetEditingNodeId(nodeId)
  }

  const handleNodeEditComplete = async (nodeId: string, content: string) => {
    // 空文字列でも更新できるようにする（空のノードを許可）
    await onNodeUpdate(nodeId, content)
    onSetEditingNodeId(null)
  }

  const handleNodeEditCancel = () => {
    onSetEditingNodeId(null)
  }

  const renderConnector = (isRoot: boolean, isLastSibling: boolean) => {
    if (isRoot) {
      return <div className="w-6 flex-shrink-0" />
    }

    return (
      <div className="relative w-6 flex-shrink-0">
        <span className="absolute left-1/2 top-0 bottom-1/2 w-px bg-gray-300 -translate-x-1/2" />
        {!isLastSibling && (
          <span className="absolute left-1/2 top-1/2 bottom-0 w-px bg-gray-300 -translate-x-1/2" />
        )}
        <svg
          className="absolute left-1/2 top-1/2 -translate-y-1/2"
          width="12"
          height="12"
          viewBox="0 0 12 12"
        >
          <path d="M0 0 V6 Q0 6 6 6 H12" stroke="#D1D5DB" strokeWidth="1" fill="none" />
        </svg>
      </div>
    )
  }

  const handleKeyDown = async (
    e: React.KeyboardEvent,
    nodeId: string,
    parentNodeId: string | null,
    currentContent: string
  ) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      if (!treeStructure.nodeMap.has(nodeId)) return

      await handleNodeEditComplete(nodeId, currentContent)

      // Get siblings to determine order_index
      const siblings = treeStructure.childrenMap.get(parentNodeId || 'root') || []
      const currentIndex = siblings.findIndex(s => s.node.id === nodeId)
      const orderIndex = currentIndex >= 0 ? currentIndex + 1 : siblings.length

      try {
        const newNode = await onNodeCreate('', parentNodeId, 'neutral', orderIndex)
        // 新しく作成したノードを選択して編集モードにする
        // 編集モードはpage.tsxのhandleNodeCreateで設定される
        onSelectNode(newNode.id)
      } catch (error) {
        console.error('Failed to create sibling node:', error)
      }
    } else if (e.key === 'Tab') {
      e.preventDefault()
      if (!treeStructure.nodeMap.has(nodeId)) return

      await handleNodeEditComplete(nodeId, currentContent)

      // Get children to determine order_index
      const children = treeStructure.childrenMap.get(nodeId) || []
      const orderIndex = children.length

      try {
        const newNode = await onNodeCreate('', nodeId, 'neutral', orderIndex)
        // 新しく作成したノードを選択して編集モードにする
        // 編集モードはpage.tsxのhandleNodeCreateで設定される
        onSelectNode(newNode.id)
      } catch (error) {
        console.error('Failed to create child node:', error)
      }
    } else if (e.key === 'Escape') {
      e.preventDefault()
      handleNodeEditCancel()
    }
  }

  const renderNode = useCallback((node: Node, parentNodeId: string | null, depth: number = 0) => {
    const children = treeStructure.childrenMap.get(node.id) || []
    const edge = edges.find(e => e.child_node_id === node.id)
    const isSelected = selectedNodeId === node.id
    const isEditing = editingNodeId === node.id
    const siblings = treeStructure.childrenMap.get(parentNodeId || 'root') || []
    const siblingIndex = siblings.findIndex(s => s.node.id === node.id)
    const isLastSibling = siblingIndex === siblings.length - 1
    const isRoot = parentNodeId === null

    return (
      <div key={node.id} className="mb-2">
        <div
          className="flex items-start gap-2"
          style={{ marginLeft: `${depth * 24}px` }}
        >
          {renderConnector(isRoot, isLastSibling)}
          {edge && (
            <span className="text-xs text-gray-500 mt-2">
              {edge.relation === 'custom' && edge.relation_label
                ? edge.relation_label
                : edge.relation === 'neutral'
                ? ''
                : edge.relation}
            </span>
          )}
          <NodeCard
            node={node}
            isSelected={isSelected}
            isEditing={isEditing}
            onClick={() => handleNodeClick(node.id)}
            onDoubleClick={() => handleNodeDoubleClick(node.id)}
            onEditComplete={(content) => handleNodeEditComplete(node.id, content)}
            onEditCancel={handleNodeEditCancel}
            onKeyDown={(e, content) => handleKeyDown(e, node.id, parentNodeId, content)}
            onDelete={() => onNodeDelete(node.id)}
          />
        </div>
        {children.map(({ node: childNode, edge: childEdge }) =>
          renderNode(childNode, node.id, depth + 1)
        )}
      </div>
    )
  }, [selectedNodeId, editingNodeId, treeStructure, edges, onNodeDelete])

  // Handle keydown for empty tree
  const handleEmptyTreeKeyDown = async (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      try {
        const newNode = await onNodeCreate('', null, 'neutral', 0)
        // 新しく作成したノードを選択して編集モードにする
        // 編集モードはpage.tsxのhandleNodeCreateで設定される
        onSelectNode(newNode.id)
      } catch (error) {
        console.error('Failed to create root node:', error)
      }
    }
  }

  // Render root nodes
  const rootNodes = treeStructure.childrenMap.get('root') || []

  return (
    <div className="p-8">
      {rootNodes.length === 0 ? (
        <div
          className="text-center text-gray-500 mt-8 cursor-pointer focus:outline-none"
          tabIndex={0}
          onKeyDown={handleEmptyTreeKeyDown}
          onClick={async () => {
            try {
              const newNode = await onNodeCreate('', null, 'neutral', 0)
              // 新しく作成したノードを選択して編集モードにする
              // 編集モードはpage.tsxのhandleNodeCreateで設定される
              onSelectNode(newNode.id)
            } catch (error) {
              console.error('Failed to create root node:', error)
            }
          }}
        >
          ノードがありません。Enterキーまたはクリックでノードを追加してください。
        </div>
      ) : (
        rootNodes.map(({ node, edge }) => renderNode(node, null, 0))
      )}
    </div>
  )
}
