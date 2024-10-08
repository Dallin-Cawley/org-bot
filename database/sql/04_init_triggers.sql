\connect org_bot

/*****************************
 * UPDATED_AT
 ****************************/
CREATE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    column_record record;
BEGIN
    FOR column_record IN
        SELECT * FROM information_schema.columns
        WHERE column_name = 'updated_at' AND table_schema = 'org_bot_schema'
    LOOP
        EXECUTE format('CREATE TRIGGER set_updated_at
                        BEFORE UPDATE ON %I.%I
                        FOR EACH ROW EXECUTE PROCEDURE set_updated_at()',
                        column_record.table_schema, column_record.table_name);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

/*****************************
 * CREATED_AT
 ****************************/
CREATE FUNCTION set_created_at() RETURNS trigger AS $$
BEGIN
    NEW.created_at := NOW();
    NEW.updated_at := Now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
    DECLARE
        column_record record;
    BEGIN
        FOR column_record IN
            SELECT * FROM information_schema.columns
            WHERE column_name = 'created_at' AND table_schema = 'org_bot_schema'
            LOOP
                EXECUTE format('CREATE TRIGGER set_created_at
                        BEFORE INSERT ON %I.%I
                        FOR EACH ROW EXECUTE PROCEDURE set_created_at()',
                               column_record.table_schema, column_record.table_name);
            END LOOP;
    END;
$$ LANGUAGE plpgsql;