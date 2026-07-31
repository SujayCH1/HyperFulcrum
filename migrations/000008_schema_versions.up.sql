CREATE TABLE schema_versions (
    id UUID PRIMARY KEY,

    project_id UUID NOT NULL,

    raw_sql TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL ,
    updated_at TIMESTAMPTZ NOT NULL 
);
