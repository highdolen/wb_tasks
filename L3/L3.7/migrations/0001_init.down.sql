DROP TRIGGER IF EXISTS items_audit_trigger ON items;
DROP FUNCTION IF EXISTS audit_item_changes();
DROP TABLE IF EXISTS item_history;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;