'use client'

import { useState, useEffect, useRef } from 'react'
import type { Node } from '@/types'

interface NodeCardProps {
  node: Node
  isSelected: boolean
  isEditing: boolean
  onClick: () => void
  onDoubleClick: () => void
  onEditComplete: (content: string) => void
  onEditCancel: () => void
  onKeyDown: (e: React.KeyboardEvent) => void
  onDelete: () => void
}

export default function NodeCard({
  node,
  isSelected,
  isEditing,
  onClick,
  onDoubleClick,
  onEditComplete,
  onEditCancel,
  onKeyDown,
  onDelete,
}: NodeCardProps) {
  const [content, setContent] = useState(node.content)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    setContent(node.content)
  }, [node.content])

  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus()
      inputRef.current.select()
    }
  }, [isEditing])

  const handleBlur = () => {
    if (content.trim() && content !== node.content) {
      onEditComplete(content)
    } else {
      setContent(node.content)
      onEditCancel()
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && e.shiftKey) {
      e.preventDefault()
      handleBlur()
      return
    }
    onKeyDown(e)
  }

  const handleContextMenu = (e: React.MouseEvent) => {
    e.preventDefault()
    if (confirm('このノードを削除しますか？')) {
      onDelete()
    }
  }

  return (
    <div
      className={`relative inline-block min-w-[200px] ${
        isSelected ? 'ring-2 ring-blue-500' : 'hover:ring-1 hover:ring-gray-300'
      } rounded-lg bg-white border border-gray-200 p-3 cursor-pointer transition-all`}
      onClick={onClick}
      onDoubleClick={onDoubleClick}
      onContextMenu={handleContextMenu}
    >
      {isEditing ? (
        <input
          ref={inputRef}
          type="text"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          onBlur={handleBlur}
          onKeyDown={handleKeyDown}
          className="w-full px-2 py-1 border border-blue-500 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          maxLength={200}
        />
      ) : (
        <div className="text-gray-900">{node.content}</div>
      )}
    </div>
  )
}
