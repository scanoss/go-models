DROP TABLE IF EXISTS db_version;
CREATE TABLE db_version
(
    package_name   text not null,
    schema_version text not null,
    created_at     text not null,
    db_release     text not null
);

INSERT INTO db_version (package_name, schema_version, created_at, db_release)
VALUES ('base', '1.0.0', '2026-01-15T10:30:00Z', '2026.01');