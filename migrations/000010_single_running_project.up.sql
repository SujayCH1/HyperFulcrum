CREATE UNIQUE INDEX projects_single_running
ON projects (running)
WHERE running = TRUE;
