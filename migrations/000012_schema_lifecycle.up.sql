ALTER TABLE schema_versions
ADD COLUMN revision BIGINT NOT NULL DEFAULT 1,
ADD COLUMN locked BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN activated_at TIMESTAMPTZ,
ADD CONSTRAINT schema_versions_revision_positive
CHECK (revision > 0);
