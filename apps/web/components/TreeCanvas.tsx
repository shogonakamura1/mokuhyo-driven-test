'use client'

import { useState, useCallback, useMemo } from 'react'
import NodeCard from './NodeCard'
import type { Node, Edge, RelationType } from '@/types'

interface TreeCanvasProps {
  nodes: Node[]
  edges: Edge[]
  selectedNodeId: string | null
  onSelectNode: (nodeId: string | null) => void
  onNodeCreate: (content: string, parentNodeId: string | null, relation?: string, orderIndex?: number) => Promise<Node>
  onNodeUpdate: (nodeId: string, content: string) => Promise<void>
  onNodeDelete: (nodeId: string) => Promise<void>
  onEdgeUpdate: (edgeId: string, relation: string, relationLabel?: string) => Promise<void>
}

export default function TreeCanvas({
  nodes,
  edges,
  selectedNodeId,
  onSelectNode,
  onNodeCreate,
  onNodeUpdate,
  onNodeDelete,
  onEdgeUpdate,
}: TreeCanvasProps) {
  const [editingNodeId, setEditingNodeId] = useState<string | null>(null)

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
    setEditingNodeId(nodeId)
  }

  const handleNodeEditComplete = async (nodeId: string, content: string) => {
    if (content.trim()) {
      await onNodeUpdate(nodeId, content)
    }
    setEditingNodeId(null)
  }

  const handleNodeEditCancel = () => {
    setEditingNodeId(null)
  }

  const handleKeyDown = async (
    e: React.KeyboardEvent,
    nodeId: string,
    parentNodeId: string | null
  ) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      const node = treeStructure.nodeMap.get(nodeId)
      if (!node) return

      await handleNodeEditComplete(nodeId, node.content)

      // Get siblings to determine order_index
      const siblings = treeStructure.childrenMap.get(parentNodeId || 'root') || []
      const currentIndex = siblings.findIndex(s => s.node.id === nodeId)
      const orderIndex = currentIndex >= 0 ? currentIndex + 1 : siblings.length

      try {
        await onNodeCreate('', parentNodeId, 'neutral', orderIndex)
      } catch (error) {
        console.error('Failed to create sibling node:', error)
      }
    } else if (e.key === 'Tab') {
      e.preventDefault()
      const node = treeStructure.nodeMap.get(nodeId)
      if (!node) return

      await handleNodeEditComplete(nodeId, node.content)

      // Get children to determine order_index
      const children = treeStructure.childrenMap.get(nodeId) || []
      const orderIndex = children.length

      try {
        await onNodeCreate('', nodeId, 'neutral', orderIndex)
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

    return (
      <div key={node.id} className="mb-2">
        <div
          className="flex items-start gap-2"
          style={{ marginLeft: `${depth * 24}px` }}
        >
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
            onKeyDown={(e) => handleKeyDown(e, node.id, parentNodeId)}
            onDelete={() => onNodeDelete(node.id)}
          />
        </div>
        {children.map(({ node: childNode, edge: childEdge }) =>
          renderNode(childNode, node.id, depth + 1)
        )}
      </div>
    )
  }, [selectedNodeId, editingNodeId, treeStructure, edges, onNodeDelete])

  // Render root nodes
  const rootNodes = treeStructure.childrenMap.get('root') || []

  return (
    <div className="p-8">
      {rootNodes.length === 0 ? (
        <div className="text-center text-gray-500 mt-8">
          ノードがありません。Enterキーでノードを追加してください。
        </div>
      ) : (
        rootNodes.map(({ node, edge }) => renderNode(node, null, 0))
      )}
    </div>
  )
}
