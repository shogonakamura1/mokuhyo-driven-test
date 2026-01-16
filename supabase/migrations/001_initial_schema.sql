-- Initial schema for goal-driven-floating-note app

-- projects
create table if not exists projects (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null,
  title text not null check (char_length(title) between 3 and 20),
  description text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  archived_at timestamptz
);
create index if not exists projects_user_id_idx on projects(user_id);

-- nodes
create table if not exists nodes (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references projects(id) on delete cascade,
  content text not null check (char_length(content) between 1 and 200),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create index if not exists nodes_project_id_idx on nodes(project_id);
create index if not exists nodes_deleted_at_idx on nodes(deleted_at) where deleted_at is null;

-- relation_type enum
do $$
begin
  if not exists (select 1 from pg_type where typname = 'relation_type') then
    create type relation_type as enum ('neutral', 'why', 'concrete', 'how', 'what', 'custom');
  end if;
end$$;

-- edges
create table if not exists edges (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references projects(id) on delete cascade,
  parent_node_id uuid references nodes(id) on delete cascade,
  child_node_id uuid not null references nodes(id) on delete cascade,
  relation relation_type not null default 'neutral',
  relation_label text check (char_length(relation_label) <= 20),
  order_index int not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint edges_unique_child unique (child_node_id)
);
create index if not exists edges_project_parent_idx on edges(project_id, parent_node_id, order_index);
create index if not exists edges_project_child_idx on edges(project_id, child_node_id);

-- user_settings
create table if not exists user_settings (
  user_id uuid primary key,
  theme text not null default 'light' check (theme in ('light', 'dark')),
  accent_color text not null default 'blue',
  updated_at timestamptz not null default now()
);

-- snapshots (optional, for future use)
create table if not exists snapshots (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references projects(id) on delete cascade,
  version int not null,
  payload jsonb not null,
  created_at timestamptz not null default now(),
  constraint snapshots_unique_version unique (project_id, version)
);

-- Row Level Security (RLS)
alter table projects enable row level security;
alter table nodes enable row level security;
alter table edges enable row level security;
alter table user_settings enable row level security;

-- projects policies
create policy "projects_select_own" on projects
for select using (auth.uid() = user_id);

create policy "projects_insert_own" on projects
for insert with check (auth.uid() = user_id);

create policy "projects_update_own" on projects
for update using (auth.uid() = user_id);

create policy "projects_delete_own" on projects
for delete using (auth.uid() = user_id);

-- nodes policies
create policy "nodes_select_own" on nodes
for select using (exists (
  select 1 from projects p where p.id = nodes.project_id and p.user_id = auth.uid()
));

create policy "nodes_insert_own" on nodes
for insert with check (exists (
  select 1 from projects p where p.id = nodes.project_id and p.user_id = auth.uid()
));

create policy "nodes_update_own" on nodes
for update using (exists (
  select 1 from projects p where p.id = nodes.project_id and p.user_id = auth.uid()
));

create policy "nodes_delete_own" on nodes
for delete using (exists (
  select 1 from projects p where p.id = nodes.project_id and p.user_id = auth.uid()
));

-- edges policies
create policy "edges_select_own" on edges
for select using (exists (
  select 1 from projects p where p.id = edges.project_id and p.user_id = auth.uid()
));

create policy "edges_insert_own" on edges
for insert with check (exists (
  select 1 from projects p where p.id = edges.project_id and p.user_id = auth.uid()
));

create policy "edges_update_own" on edges
for update using (exists (
  select 1 from projects p where p.id = edges.project_id and p.user_id = auth.uid()
));

create policy "edges_delete_own" on edges
for delete using (exists (
  select 1 from projects p where p.id = edges.project_id and p.user_id = auth.uid()
));

-- user_settings policies
create policy "settings_select_own" on user_settings
for select using (auth.uid() = user_id);

create policy "settings_insert_own" on user_settings
for insert with check (auth.uid() = user_id);

create policy "settings_update_own" on user_settings
for update using (auth.uid() = user_id);
