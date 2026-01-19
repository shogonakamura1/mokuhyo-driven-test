'use client'

import { useState, useEffect, useRef } from 'react'
import type { Node } from '@/types'

interface NodeCardProps {
  node: Node
  isSelected: boolean
  isEditing: boolean
  onClick: () => void
  onDoubleClick: () => void
  onEditComplete: (content: string) => Promise<void>
  onEditCancel: () => void
  onKeyDown: (e: React.KeyboardEvent, content: string) => void
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

  const handleBlur = async () => {
    // 空文字列でも更新できるようにする（空のノードを許可）
    if (content !== node.content) {
      await onEditComplete(content)
    } else {
      setContent(node.content)
      onEditCancel()
    }
  }

  const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.nativeEvent.isComposing) {
      return
    }
    if (e.key === 'Enter') {
      e.preventDefault()
      // Shift+Enterは保存して編集終了、EnterはTreeCanvas側で保存とノード作成を行う
      if (e.shiftKey) {
        if (content !== node.content) {
          await onEditComplete(content)
        }
        return
      }
      onKeyDown(e, content)
      return
    }
    onKeyDown(e, content)
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
        <div className={`text-gray-900 ${node.content === '' ? 'text-gray-400' : ''}`}>
          {node.content || '\u00A0'}
        </div>
      )}
    </div>
  )
}
