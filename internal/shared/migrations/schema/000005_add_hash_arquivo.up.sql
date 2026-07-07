ALTER TABLE arquivo_importado
  ADD COLUMN IF NOT EXISTS hash_sha256 VARCHAR(64) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_arquivo_importado_hash
  ON arquivo_importado (hash_sha256);
