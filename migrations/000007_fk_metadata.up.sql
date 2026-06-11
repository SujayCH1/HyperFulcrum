CREATE TABLE fk_edges (
    project_id UUID NOT NULL,

    parent_table TEXT NOT NULL,
    parent_column TEXT NOT NULL,

    child_table TEXT NOT NULL,
    child_column TEXT NOT NULL,

    PRIMARY KEY (
        project_id,
        parent_table,
        parent_column,
        child_table,
        child_column
    )
);

CREATE INDEX idx_fk_edges_project_parent
    ON fk_edges(project_id, parent_table, parent_column);

CREATE INDEX idx_fk_edges_project_child
    ON fk_edges(project_id, child_table, child_column);
