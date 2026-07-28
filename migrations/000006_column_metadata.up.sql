CREATE TABLE IF NOT EXISTS columns (
    project_id UUID NOT NULL,
    table_name TEXT NOT NULL,
    column_name TEXT NOT NULL,

    data_type TEXT NOT NULL,
    nullable BOOLEAN NOT NULL,
    is_primary_key BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (project_id, table_name, column_name)
);

CREATE INDEX idx_columns_project_table
    ON columns(project_id, table_name);
