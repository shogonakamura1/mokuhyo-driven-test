-- Add users table for Google authentication

create table if not exists users (
  id uuid primary key default gen_random_uuid(),
  google_id text not null unique,
  email text not null,
  name text not null,
  picture text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists users_google_id_idx on users(google_id);
create index if not exists users_email_idx on users(email);

-- Update projects table to reference users table
-- Note: This assumes projects.user_id already exists and references auth.users in Supabase
-- For local PostgreSQL, we'll need to ensure the foreign key relationship
alter table projects 
  add constraint projects_user_id_fkey 
  foreign key (user_id) references users(id) on delete cascade;

-- Update user_settings table to reference users table
alter table user_settings
  add constraint user_settings_user_id_fkey
  foreign key (user_id) references users(id) on delete cascade;
