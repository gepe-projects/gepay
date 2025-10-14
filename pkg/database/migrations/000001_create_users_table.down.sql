-- Down migration
DROP INDEX IF EXISTS name;
DROP INDEX IF EXISTS idx_user_identities_provider_provider_id;

DROP TABLE IF EXISTS user_identities;
DROP TABLE IF EXISTS users;