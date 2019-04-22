package storage

const schemaFmt = `-- listed_tokens table
CREATE TABLE IF NOT EXISTS "%[1]s"
(
    id        SERIAL PRIMARY KEY,
    address   text      NOT NULL UNIQUE,
    name      text      NOT NULL,
    symbol    text      NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    parent_id INT REFERENCES "%[1]s" (id)
);

-- version table
CREATE TABLE IF NOT EXISTS "%[2]s"
(
    id           SERIAL PRIMARY KEY,
    version      INT    NOT NULL,
    block_number bigint NOT NULL
);

INSERT INTO "%[2]s"(version, block_number)
SELECT 0, 0
WHERE (SELECT COUNT(*) FROM "%[2]s") = 0;


-- reserves table
CREATE TABLE IF NOT EXISTS "%[3]s"
(
    id      serial NOT NULL PRIMARY KEY,
    address TEXT   NOT NULL UNIQUE
);

-- reserve_token table
CREATE TABLE IF NOT EXISTS "%[4]s"
(
    id         SERIAL NOT NULL,
    token_id   INT REFERENCES "%[1]s" (id),
    reserve_id INT REFERENCES "%[3]s" (id),
    PRIMARY KEY (token_id, reserve_id)
);

-- save_token function saves or update given token to database and return TRUE if anything changes recorded to database.
CREATE OR REPLACE FUNCTION save_token(_address "%[1]s".address%%TYPE,
                                      _name "%[1]s".name%%TYPE,
                                      _symbol "%[1]s".name%%TYPE,
                                      _timestamp "%[1]s".timestamp%%TYPE,
                                      _parent_address "%[1]s".address%%TYPE,
                                      _reserve_address "%[3]s".address%%TYPE) RETURNS boolean AS
$$
DECLARE
    stored_token   RECORD;
    stored_reserve RECORD;
    changed        BOOLEAN = FALSE;
    stored_parent_id      INTEGER = NULL;
BEGIN
    IF _parent_address IS NOT NULl THEN
        SELECT id INTO stored_parent_id FROM "%[1]s" WHERE address = _parent_address;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'token with id %% does not exist', stored_parent_id;
        END IF;
    END IF;
    
    SELECT * INTO
        stored_token
    FROM "%[1]s"
    WHERE address = _address;
    IF NOT FOUND THEN
        INSERT INTO "%[1]s"(address, name, symbol, timestamp, parent_id)
        VALUES (_address, _name, _symbol, _timestamp, stored_parent_id) RETURNING * INTO stored_token;
        changed := TRUE;
    ELSE
        IF FOUND THEN
            IF stored_token.address != _address OR
               stored_token.name != _name OR
               stored_token.symbol != _symbol OR
               stored_token.timestamp != _timestamp OR
               (stored_token.parent_id IS NULL AND stored_parent_id IS NOT NULL) OR
               (stored_token.parent_id IS NOT NULL AND stored_parent_id IS NULL) OR
               stored_token.parent_id != stored_parent_id THEN
                UPDATE "%[1]s"
                SET name=_name,
                    symbol=_symbol,
                    timestamp=_timestamp,
                    parent_id = stored_parent_id
                WHERE address = _address RETURNING * INTO stored_token;
                changed := TRUE;
            END IF;
        END IF;
    END IF;


    INSERT INTO "%[3]s"(address)
    VALUES (_reserve_address)
    ON CONFLICT DO NOTHING;
    SELECT * INTO stored_reserve FROM "%[3]s" WHERE address = _reserve_address;

    INSERT INTO "%[4]s"(token_id, reserve_id)
    VALUES (stored_token.id, stored_reserve.id)
    ON CONFLICT DO NOTHING;
    IF FOUND THEN
        changed := TRUE;
    END IF;

    RETURN changed;
END
$$ LANGUAGE PLPGSQL;

CREATE OR REPLACE VIEW tokens_view AS
SELECT joined.address as address,
       joined.name,
       joined.symbol,
       joined.timestamp,
       joined.reserve_address,
       array_agg(joined.old_address)
                 FILTER ( WHERE joined.old_address IS NOT NULL)::text[]     AS old_addresses,
       array_agg(extract(EPOCH FROM joined.old_timestamp) * 1000)
                 FILTER ( WHERE joined.old_timestamp IS NOT NULL)::BIGINT[] AS old_timestamps
FROM (SELECT toks.address,
             toks.name,
             toks.symbol,
             toks.timestamp,
             reserve.address as reserve_address,
             olds.address   AS old_address,
             olds.timestamp AS old_timestamp
      FROM "%[1]s" AS toks
               LEFT JOIN "%[1]s" AS olds
                         ON toks.id = olds.parent_id
               JOIN "%[4]s" as token_reserve 
                         ON toks.id = token_reserve.token_id
               JOIN "%[3]s" as reserve
                         ON reserve.id = token_reserve.reserve_id
      WHERE toks.parent_id IS NULL
      ORDER BY timestamp DESC) AS joined
GROUP BY joined.address, joined.name, joined.symbol, joined.timestamp, joined.reserve_address;
`
