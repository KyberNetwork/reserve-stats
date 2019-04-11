package storage

const schemaFmt = `CREATE TABLE IF NOT EXISTS "app_names"
(
    id     SERIAL PRIMARY KEY,
    name   text NOT NULL UNIQUE CHECK ( LENGTH(name) > 0 ),
    active boolean DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS "addresses"
(
    id          SERIAL PRIMARY KEY,
    address     text   NOT NULL UNIQUE CHECK ( LENGTH(address) > 0 ),
    app_name_id SERIAL NOT NULL REFERENCES app_names (id)
);

-- create_or_update_app creates or update the application with given name.
-- If id is provided, this function will only update the application with given id if exists.
-- The list of addresses will be create and link to the created/updated application. Other existing addresses that
-- currently link to application will be removed.
CREATE OR REPLACE FUNCTION create_or_update_app(INOUT _id app_names.id%TYPE,
                                                _name app_names.name%TYPE,
                                                _addresses TEXT[],
                                                OUT updated BOOLEAN) AS
$$
DECLARE
    _address addresses.address%TYPE;
BEGIN
    updated = FALSE;
    IF _id = 0 THEN
        INSERT INTO app_names (name, active) VALUES (_name, TRUE) ON CONFLICT DO NOTHING RETURNING id INTO _id;
        IF _id IS NULL THEN
            updated = TRUE;
            UPDATE app_names SET active = TRUE WHERE name = _name RETURNING id INTO _id;
        END IF;
    ELSE
        updated = TRUE;
        UPDATE app_names SET name = _name, active = TRUE WHERE id = _id;
    END IF;


    IF _addresses IS NOT NULL THEN
        FOREACH _address IN ARRAY _addresses
            LOOP
                INSERT INTO "addresses"(address, app_name_id)
                VALUES (_address, _id)
                ON CONFLICT ON CONSTRAINT addresses_address_key DO UPDATE SET app_name_id = EXCLUDED.app_name_id;
            END LOOP;

        DELETE
        FROM addresses
        WHERE app_name_id = _id
          AND NOT address = ANY (_addresses);
    END IF;

    RETURN;
END;
$$ LANGUAGE PLPGSQL;

-- update_app update app with given id with given name and addresses if they are not null and zero.
CREATE OR REPLACE FUNCTION update_app(_id app_names.id%TYPE, _name app_names.name%TYPE,
                                      _addresses TEXT[]) RETURNS VOID AS
$$
DECLARE
    _address addresses.address%TYPE;
BEGIN
    PERFORM id FROM app_names WHERE id = _id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'application with id % does not exists', _id USING ERRCODE = 'no_data_found';
    END IF;

    UPDATE app_names SET active = TRUE WHERE id = _id;

    IF _name IS NOT NULL AND LENGTH(_name) <> 0 THEN
        UPDATE app_names SET name = _name WHERE id = _id;
    END IF;

    IF _addresses IS NOT NULL THEN
        IF _addresses IS NOT NULL AND ARRAY_LENGTH(_addresses, 1) <> 0 THEN
            FOREACH _address IN ARRAY _addresses
                LOOP
                    INSERT INTO "addresses"(address, app_name_id)
                    VALUES (_address, _id)
                    ON CONFLICT ON CONSTRAINT addresses_address_key DO UPDATE SET app_name_id = EXCLUDED.app_name_id;
                END LOOP;

            DELETE
            FROM addresses
            WHERE app_name_id = _id
              AND NOT address = ANY (_addresses);
        END IF;
    END IF;
    RETURN;
END
$$ LANGUAGE PLPGSQL
`
