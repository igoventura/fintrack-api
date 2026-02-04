-- Update transactions schema

ALTER TABLE transactions
ADD COLUMN currency VARCHAR(3) NOT NULL,
ADD COLUMN parent_transaction_id UUID,
DROP COLUMN previous_sibling_transaction_id,
DROP COLUMN next_sibling_transaction_id;

ALTER TABLE transactions ADD FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id);

ALTER TABLE transactions_tags
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN deactivated_at TIMESTAMPTZ;

---- create above / drop below ----

ALTER TABLE transactions_tags
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deactivated_at;

ALTER TABLE transactions
DROP CONSTRAINT transactions_parent_transaction_id_fkey;

ALTER TABLE transactions
DROP COLUMN currency,
DROP COLUMN parent_transaction_id,
ADD COLUMN previous_sibling_transaction_id UUID,
ADD COLUMN next_sibling_transaction_id UUID;

ALTER TABLE transactions ADD FOREIGN KEY (previous_sibling_transaction_id) REFERENCES transactions(id);
ALTER TABLE transactions ADD FOREIGN KEY (next_sibling_transaction_id) REFERENCES transactions(id);
