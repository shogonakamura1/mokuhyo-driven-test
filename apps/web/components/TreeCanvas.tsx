'use client' // Next.jsのクライアントコンポーネントであることを示す

import { useCallback, useMemo } from 'react' // Reactのフックをインポート
import NodeCard from './NodeCard' // ノードカードコンポーネントをインポート
import type { Node, Edge, RelationType } from '@/types' // 型定義をインポート

interface TreeCanvasProps { // TreeCanvasコンポーネントのプロップスの型定義
  nodes: Node[] // ノードの配列
  edges: Edge[] // エッジ（ノード間の関係）の配列
  selectedNodeId: string | null // 選択されているノードのID
  editingNodeId: string | null // 編集中のノードのID
  onSelectNode: (nodeId: string | null) => void // ノード選択時のコールバック
  onSetEditingNodeId: (nodeId: string | null) => void // 編集ノードID設定時のコールバック
  onNodeCreate: (
    content: string,
    parentNodeId: string | null,
    relation?: string,
    orderIndex?: number,
    question?: string
  ) => Promise<Node> // ノード作成時のコールバック
  onNodeUpdate: (nodeId: string, content: string) => Promise<void> // ノード更新時のコールバック
  onNodeDelete: (nodeId: string) => Promise<void> // ノード削除時のコールバック
  onEdgeUpdate: (edgeId: string, relation: string, relationLabel?: string) => Promise<void> // エッジ更新時のコールバック
}

export default function TreeCanvas({ // TreeCanvasコンポーネントのエクスポート
  nodes, // ノードの配列を受け取る
  edges, // エッジの配列を受け取る
  selectedNodeId, // 選択されているノードIDを受け取る
  editingNodeId, // 編集中のノードIDを受け取る
  onSelectNode, // ノード選択コールバックを受け取る
  onSetEditingNodeId, // 編集ノードID設定コールバックを受け取る
  onNodeCreate, // ノード作成コールバックを受け取る
  onNodeUpdate, // ノード更新コールバックを受け取る
  onNodeDelete, // ノード削除コールバックを受け取る
  onEdgeUpdate, // エッジ更新コールバックを受け取る
}: TreeCanvasProps) {

  // Build tree structure
  const treeStructure = useMemo(() => { // ツリー構造をメモ化（nodes/edgesが変わった時のみ再計算）
    const nodeMap = new Map(nodes.map(node => [node.id, node])) // ノードIDをキーとしたMapを作成
    const childrenMap = new Map<string, Array<{ node: Node; edge: Edge }>>() // 親ノードIDをキーとした子ノード配列のMapを作成

    edges.forEach(edge => { // 各エッジを処理
      const parentId = edge.parent_node_id || 'root' // 親ノードIDを取得（nullの場合は'root'）
      if (!childrenMap.has(parentId)) { // 親ノードIDがMapに存在しない場合
        childrenMap.set(parentId, []) // 空配列を設定
      }
      const node = nodeMap.get(edge.child_node_id) // 子ノードを取得
      if (node) { // ノードが存在する場合
        childrenMap.get(parentId)!.push({ node, edge }) // 子ノード配列に追加
      }
    })

    // Sort children by order_index
    childrenMap.forEach(children => { // 各親ノードの子ノード配列を処理
      children.sort((a, b) => a.edge.order_index - b.edge.order_index) // order_indexでソート
    })

    return { nodeMap, childrenMap } // ノードMapと子ノードMapを返す
  }, [nodes, edges]) // nodesとedgesが変わった時のみ再計算

  const handleNodeClick = (nodeId: string) => { // ノードクリック時のハンドラ
    onSelectNode(nodeId) // ノード選択コールバックを呼ぶ
    onSetEditingNodeId(nodeId) // クリックで編集モードにする
  }

  const handleNodeDoubleClick = (nodeId: string) => { // ノードダブルクリック時のハンドラ
    onSetEditingNodeId(nodeId) // 編集ノードID設定コールバックを呼ぶ
  }

  const handleNodeEditComplete = async (nodeId: string, content: string) => { // ノード編集完了時のハンドラ
    // 空文字列でも更新できるようにする（空のノードを許可）
    await onNodeUpdate(nodeId, content) // ノード更新コールバックを呼ぶ
    onSetEditingNodeId(null) // 編集モードを解除
  }

  const handleNodeEditCancel = () => { // ノード編集キャンセル時のハンドラ
    onSetEditingNodeId(null) // 編集モードを解除
  }

  const renderChildConnector = (isLastSibling: boolean) => ( // 子ノードへの接続線を描画する関数
    <div className="relative w-6 flex-shrink-0"> {/* 相対配置で幅6のコンテナ */}
      <span className="absolute left-1/2 top-0 bottom-1/2 w-px bg-gray-300 -translate-x-1/2" /> {/* 上から中央までの縦線 */}
      {/* 最後の兄弟ノードでない場合、中央から下までの縦線を表示 */}
      {!isLastSibling && (
        <span className="absolute left-1/2 top-1/2 bottom-0 w-px bg-gray-300 -translate-x-1/2" />
      )}
      <svg // SVG要素で曲線を描画
        className="absolute left-1/2 top-1/2 -translate-y-1/2" // 中央に配置
        width="12" // SVGの幅
        height="12" // SVGの高さ
        viewBox="0 0 12 12" // ビューボックス
      >
        <path d="M0 0 V6 Q0 6 6 6 H12" stroke="#D1D5DB" strokeWidth="1" fill="none" /> {/* 縦線→曲線→横線のパス */}
      </svg>
    </div>
  )

  const handleKeyDown = async ( // キーボード入力時のハンドラ
    e: React.KeyboardEvent, // キーボードイベント
    nodeId: string, // 現在のノードID
    parentNodeId: string | null, // 親ノードID
    currentContent: string // 現在の入力内容
  ) => {
    if (e.nativeEvent.isComposing) {
      return
    }
    if (e.key === 'Enter' && !e.shiftKey) { // Enterキー（Shiftなし）が押された場合
      e.preventDefault() // デフォルト動作をキャンセル
      if (!treeStructure.nodeMap.has(nodeId)) return // ノードが存在しない場合は終了

      await handleNodeEditComplete(nodeId, currentContent) // 編集内容を保存

      const currentNode = treeStructure.nodeMap.get(nodeId)
      const inheritedQuestion = currentNode?.question?.trim()

      // Get siblings to determine order_index
      const siblings = treeStructure.childrenMap.get(parentNodeId || 'root') || [] // 兄弟ノードを取得
      const currentIndex = siblings.findIndex(s => s.node.id === nodeId) // 現在のノードのインデックスを取得
      const orderIndex = currentIndex >= 0 ? currentIndex + 1 : siblings.length // 新しいノードのorder_indexを計算

      try {
        const newNode = await onNodeCreate(
          '',
          parentNodeId,
          'neutral',
          orderIndex,
          inheritedQuestion || undefined
        ) // 並列ノード（兄弟）を作成
        // 新しく作成したノードを選択して編集モードにする
        // 編集モードはpage.tsxのhandleNodeCreateで設定される
        onSelectNode(newNode.id) // 新しいノードを選択
      } catch (error) {
        console.error('Failed to create sibling node:', error) // エラーをログ出力
      }
    } else if (e.key === 'Tab') { // Tabキーが押された場合
      e.preventDefault() // デフォルト動作をキャンセル
      if (!treeStructure.nodeMap.has(nodeId)) return // ノードが存在しない場合は終了

      await handleNodeEditComplete(nodeId, currentContent) // 編集内容を保存

      // Get children to determine order_index
      const children = treeStructure.childrenMap.get(nodeId) || [] // 子ノードを取得
      const orderIndex = children.length // 新しいノードのorder_indexを計算

      try {
        const newNode = await onNodeCreate('', nodeId, 'neutral', orderIndex) // 次ノード（子）を作成
        // 新しく作成したノードを選択して編集モードにする
        // 編集モードはpage.tsxのhandleNodeCreateで設定される
        onSelectNode(newNode.id) // 新しいノードを選択
      } catch (error) {
        console.error('Failed to create child node:', error) // エラーをログ出力
      }
    } else if (e.key === 'Escape') { // Escapeキーが押された場合
      e.preventDefault() // デフォルト動作をキャンセル
      handleNodeEditCancel() // 編集をキャンセル
    }
  }

  const renderNode = useCallback((node: Node, parentNodeId: string | null) => { // ノードを再帰的に描画する関数（メモ化）
    const children = treeStructure.childrenMap.get(node.id) || [] // 子ノードを取得
    const edge = edges.find(e => e.child_node_id === node.id) // このノードへのエッジを取得
    const questionLabel = node.question?.trim() // AIが選んだ質問ラベル
    const edgeLabel = questionLabel
      ? questionLabel
      : edge?.relation === 'custom' && edge.relation_label
      ? edge.relation_label
      : edge?.relation === 'neutral'
      ? ''
      : edge?.relation // AI質問がない場合は既存の関係ラベルを使う
    const isSelected = selectedNodeId === node.id // 選択されているかどうか
    const isEditing = editingNodeId === node.id // 編集中かどうか
    const shouldReserveLabel = parentNodeId !== null && !!edge // 兄弟間でラベル位置を揃える

    return (
      <div key={node.id} className="flex items-start gap-6 mb-2"> {/* 横並びレイアウト、上揃え、間隔6、下マージン2 */}
        <div className="flex items-center gap-2"> {/* ノードとラベルを横並び、中央揃え、間隔2 */}
          {/* エッジが存在する場合、関係ラベルを表示 */}
          {shouldReserveLabel && (
            <span className="text-xs text-gray-500 mt-2 min-w-[7rem] max-w-[10rem] text-right"> {/* 関係ラベルを表示 */}
              {/* カスタム関係でラベルがある場合はカスタムラベル、neutral関係の場合は何も表示せず、それ以外は関係タイプを表示 */}
              {edgeLabel || ''}
            </span>
          )}
          <NodeCard // ノードカードコンポーネント
            node={node} // ノードデータを渡す
            isSelected={isSelected} // 選択状態を渡す
            isEditing={isEditing} // 編集状態を渡す
            onClick={() => handleNodeClick(node.id)} // クリックハンドラを渡す
            onDoubleClick={() => handleNodeDoubleClick(node.id)} // ダブルクリックハンドラを渡す
            onEditComplete={(content) => handleNodeEditComplete(node.id, content)} // 編集完了ハンドラを渡す
            onEditCancel={handleNodeEditCancel} // 編集キャンセルハンドラを渡す
            onKeyDown={(e, content) => handleKeyDown(e, node.id, parentNodeId, content)} // キーダウンハンドラを渡す
            onDelete={() => onNodeDelete(node.id)} // 削除ハンドラを渡す
          />
        </div>
        {/* 子ノードが存在する場合 */}
        {children.length > 0 && (
          <div className="flex flex-col gap-2"> {/* 子ノードを縦並び、間隔2 */}
            {/* 各子ノードを処理 */}
            {children.map(({ node: childNode }, index) => {
              const isLastSibling = index === children.length - 1 // 最後の兄弟かどうか
              return (
                <div key={childNode.id} className="flex items-center gap-2"> {/* 接続線とノードを横並び、中央揃え、間隔2 */}
                  {renderChildConnector(isLastSibling)} {/* 接続線を描画 */}
                  {renderNode(childNode, node.id)} {/* 子ノードを再帰的に描画 */}
                </div>
              )
            })}
          </div>
        )}
      </div>
    )
  }, [selectedNodeId, editingNodeId, treeStructure, edges, onNodeDelete]) // 依存配列：これらの値が変わった時のみ再描画

  // Handle keydown for empty tree
  const handleEmptyTreeKeyDown = async (e: React.KeyboardEvent) => { // 空のツリーでのキーダウンハンドラ
    if (e.key === 'Enter' && !e.shiftKey) { // Enterキー（Shiftなし）が押された場合
      e.preventDefault() // デフォルト動作をキャンセル
      try {
        const newNode = await onNodeCreate('', null, 'neutral', 0) // ルートノードを作成
        // 新しく作成したノードを選択して編集モードにする
        // 編集モードはpage.tsxのhandleNodeCreateで設定される
        onSelectNode(newNode.id) // 新しいノードを選択
      } catch (error) {
        console.error('Failed to create root node:', error) // エラーをログ出力
      }
    }
  }

  // Render root nodes
  const rootNodes = treeStructure.childrenMap.get('root') || [] // ルートノードを取得

  return (
    <div className="p-8"> {/* パディング8のコンテナ */}
      {rootNodes.length === 0 ? ( // ルートノードが存在しない場合
        <div
          className="text-center text-gray-500 mt-8 cursor-pointer focus:outline-none" // 中央揃え、グレー、カーソルポインター、フォーカスアウトラインなし
          tabIndex={0} // フォーカス可能にする
          onKeyDown={handleEmptyTreeKeyDown} // キーダウンハンドラを設定
          onClick={async () => { // クリックハンドラを設定
            try {
              const newNode = await onNodeCreate('', null, 'neutral', 0) // ルートノードを作成
              // 新しく作成したノードを選択して編集モードにする
              // 編集モードはpage.tsxのhandleNodeCreateで設定される
              onSelectNode(newNode.id) // 新しいノードを選択
            } catch (error) {
              console.error('Failed to create root node:', error) // エラーをログ出力
            }
          }}
        >
          ノードがありません。Enterキーまたはクリックでノードを追加してください。 {/* 空のツリー時のメッセージ */}
        </div>
      ) : ( // ルートノードが存在する場合
        rootNodes.map(({ node }) => renderNode(node, null)) // 各ルートノードを描画
      )}
    </div>
  )
}
