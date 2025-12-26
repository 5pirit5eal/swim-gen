-- Drill Embeddings - separate table for drill vectors
create table if not exists drill_embeddings (
  collection_id uuid,
  embedding extensions.vector(768),
  document varchar,
  cmetadata json,
  uuid uuid not null,
  constraint drill_embeddings_collection_id_fkey
    foreign key (collection_id) references embedders (uuid) on delete cascade,
  primary key (uuid)
);

-- Index for drill embeddings collection lookup
create index if not exists drill_embeddings_collection_id
  on drill_embeddings (collection_id);
