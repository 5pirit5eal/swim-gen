-- Embedders (Collections) - required by langchaingo pgvector
create table if not exists embedders (
  name varchar,
  cmetadata json,
  uuid uuid not null,
  unique (name),
  primary key (uuid)
);

-- Embeddings - required by langchaingo pgvector
-- Table name matches EMBEDDING_NAME from .env (embeddings)
create table if not exists embeddings (
  collection_id uuid,
  embedding vector(768),  -- dimensions from EMBEDDING_SIZE=768
  document varchar,
  cmetadata json,
  uuid uuid not null,
  constraint langchain_pg_embedding_collection_id_fkey
    foreign key (collection_id) references embedders (uuid) on delete cascade,
  primary key (uuid)
);

-- Index for embeddings collection lookup
create index if not exists embeddings_collection_id
  on embeddings (collection_id);

-- Optional: HNSW index for vector similarity search (uncomment if needed)
-- Note: This should be created AFTER inserting enough rows for good performance
-- create index if not exists embeddings_embedding_hnsw
--   on embeddings using hnsw (embedding vector_cosine_ops);

-- Plans - Gold standard plans
create table if not exists plans (
  plan_id uuid primary key DEFAULT gen_random_uuid(),
  title text not null,
  description text not null,
  plan_table jsonb not null,
  created_at timestamptz default now()
);

-- Scraped mapping
create table if not exists scraped (
  url text not null,
  collection_id uuid not null,
  plan_id uuid not null references plans(plan_id) on delete cascade,
  created_at timestamptz default now(),
  primary key (url, collection_id),
  FOREIGN KEY (collection_id) REFERENCES embedders (uuid) ON DELETE CASCADE
);

-- Donations
create table if not exists donations (
  user_id uuid not null,
  plan_id uuid not null references plans(plan_id) on delete cascade,
  created_at timestamptz default now(),
  primary key (user_id, plan_id)
);

-- Feedback
create table if not exists feedback (
  user_id uuid not null,
  plan_id uuid not null references plans(plan_id) on delete cascade,
  rating int not null check (rating between 1 and 5),
  comment text not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  primary key (user_id, plan_id)
);

-- JSONB potential future query performance
create index if not exists idx_plans_plan_table_jsonb on plans using gin (plan_table jsonb_path_ops);

-- Foreign key access patterns
create index if not exists idx_scraped_plan_id on scraped(plan_id);
create index if not exists idx_donations_plan_id on donations(plan_id);
create index if not exists idx_feedback_plan_id on feedback(plan_id);

-- Vector index (run AFTER enough rows present)
-- Example (adjust table/column names per pgvector store):
-- create index if not exists emb_idx on embeddings using ivfflat (embedding vector_cosine_ops) with (lists=100);
-- analyze embeddings;

-- Grant ownership of tables to coach user
ALTER TABLE embedders OWNER TO coach;
ALTER TABLE embeddings OWNER TO coach;
ALTER TABLE plans OWNER TO coach;
ALTER TABLE scraped OWNER TO coach;
ALTER TABLE donations OWNER TO coach;
ALTER TABLE feedback OWNER TO coach;
