export type Project = {
  id: string
  user_id: string
  title: string
  description?: string
  created_at: string
  updated_at: string
  archived_at?: string
}

export type Node = {
  id: string
  project_id: string
  content: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export type RelationType = 'neutral' | 'why' | 'concrete' | 'how' | 'what' | 'custom'

export type Edge = {
  id: string
  project_id: string
  parent_node_id?: string
  child_node_id: string
  relation: RelationType
  relation_label?: string
  order_index: number
  created_at: string
  updated_at: string
}

export type TreeResponse = {
  project: Project
  nodes: Node[]
  edges: Edge[]
}

export type UserSettings = {
  user_id: string
  theme: 'light' | 'dark'
  accent_color: string
  updated_at: string
}

export type MeResponse = {
  user: {
    id: string
  }
  settings: UserSettings
}
