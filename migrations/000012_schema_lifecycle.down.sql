ALTER TABLE schema_versions
DROP CONSTRAINT IF EXISTS schema_versions_revision_positive,
DROP COLUMN IF EXISTS activated_at,
DROP COLUMN IF EXISTS locked,
DROP COLUMN IF EXISTS revision;
