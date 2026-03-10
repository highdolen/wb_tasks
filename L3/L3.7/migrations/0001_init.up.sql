-- users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL
);

-- items
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    updated_at TIMESTAMP DEFAULT now()
);

-- item history
CREATE TABLE item_history (
    id SERIAL PRIMARY KEY,
    item_id INT,
    action TEXT,
    old_value JSONB,
    new_value JSONB,
    changed_by TEXT,
    changed_at TIMESTAMP DEFAULT now()
);

-- trigger function
DROP FUNCTION IF EXISTS audit_item_changes();

CREATE OR REPLACE FUNCTION audit_item_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO item_history (
            item_id,
            action,
            new_value,
            changed_by,
            changed_at
        )
        VALUES (
            NEW.id,
            'INSERT',
            row_to_json(NEW),
            current_setting('app.user', true),
            NOW()
        );
        RETURN NEW;

    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO item_history (
            item_id,
            action,
            old_value,
            new_value,
            changed_by,
            changed_at
        )
        VALUES (
            NEW.id,
            'UPDATE',
            row_to_json(OLD),
            row_to_json(NEW),
            current_setting('app.user', true),
            NOW()
        );
        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO item_history (
            item_id,
            action,
            old_value,
            changed_by,
            changed_at
        )
        VALUES (
            OLD.id,
            'DELETE',
            row_to_json(OLD),
            current_setting('app.user', true),
            NOW()
        );
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

--TRIGGER
CREATE TRIGGER items_audit_trigger
AFTER INSERT OR UPDATE OR DELETE
ON items
FOR EACH ROW
EXECUTE FUNCTION audit_item_changes();