CREATE TABLE IF NOT EXISTS db_schema_version(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version TEXT NOT NULL,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_db_schema_version
AFTER UPDATE ON db_schema_version
BEGIN
    UPDATE db_schema_version SET modified_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS environments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    env_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_environments
AFTER UPDATE ON environments
BEGIN
    UPDATE environments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Create the first environment - but only if it hasn't ever been created yet
INSERT INTO environments (name, env_type) SELECT 'local', 'local'
WHERE NOT EXISTS (SELECT 1 FROM db_schema_version);

-- Finally, insert the database version
INSERT INTO db_schema_version (version) SELECT '0'
WHERE NOT EXISTS (SELECT 1 FROM db_schema_version);