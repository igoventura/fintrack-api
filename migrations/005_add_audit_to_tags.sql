-- Add audit columns to tags



ALTER TABLE tags
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN created_by UUID NOT NULL,
ADD COLUMN updated_by UUID NOT NULL,
ADD COLUMN deactivated_by UUID;

ALTER TABLE tags ADD FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE tags ADD FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE tags ADD FOREIGN KEY (deactivated_by) REFERENCES users(id);

---- create above / drop below ----

ALTER TABLE tags
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN created_by,
DROP COLUMN updated_by,
DROP COLUMN deactivated_by;
