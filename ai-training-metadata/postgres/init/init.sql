CREATE TABLE processes (
    id UUID PRIMARY KEY,
    start_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    end_time TIMESTAMP WITHOUT TIME ZONE,
    project TEXT NOT NULL,
    run TEXT NOT NULL,
    config JSONB
);

CREATE TABLE processes_audit (
    id SERIAL PRIMARY KEY,
    process_id UUID NOT NULL,
    previous_state JSONB NOT NULL,
    current_state JSONB NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION processes_audit_trigger()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO processes_audit (process_id, previous_state, current_state)
    VALUES (NEW.id, OLD::JSONB, NEW::JSONB);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER processes_audit
AFTER UPDATE ON processes
FOR EACH ROW
EXECUTE FUNCTION processes_audit_trigger();
