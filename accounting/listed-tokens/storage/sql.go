package storage

const schemaFmt = `-- listed_tokens table
CREATE TABLE IF NOT EXISTS "listed_tokens"
(
    id        SERIAL PRIMARY KEY,
    address   text      NOT NULL UNIQUE,
    name      text      NOT NULL,
    symbol    text      NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    parent_id INT REFERENCES "listed_tokens" (id)
);

-- version table
CREATE TABLE IF NOT EXISTS "listed_tokens_version"
(
    id           SERIAL PRIMARY KEY,
    version      INT    NOT NULL,
    block_number bigint NOT NULL
);

INSERT INTO "listed_tokens_version"(version, block_number)
SELECT 0, 0
WHERE (SELECT COUNT(*) FROM "listed_tokens_version") = 0;


-- reserves table
CREATE TABLE IF NOT EXISTS "listed_tokens_reserves"
(
    id      serial NOT NULL PRIMARY KEY,
    address TEXT   NOT NULL UNIQUE
);

-- reserve_token table
CREATE TABLE IF NOT EXISTS "listed_tokens_reserves_tokens"
(
    id         SERIAL NOT NULL,
    token_id   INT REFERENCES "listed_tokens" (id),
    reserve_id INT REFERENCES "listed_tokens_reserves" (id),
    PRIMARY KEY (token_id, reserve_id)
);

-- save_token function saves or update given token to database and return TRUE if anything changes recorded to database.
CREATE OR REPLACE FUNCTION save_token(_address "listed_tokens".address%TYPE,
                                      _name "listed_tokens".name%TYPE,
                                      _symbol "listed_tokens".name%TYPE,
                                      _timestamp "listed_tokens".timestamp%TYPE,
                                      _parent_address "listed_tokens".address%TYPE,
                                      _reserve_address "listed_tokens_reserves".address%TYPE) RETURNS boolean AS
$$
DECLARE
    stored_token   RECORD;
    stored_reserve RECORD;
    changed        BOOLEAN = FALSE;
    stored_parent_id      INTEGER = NULL;
BEGIN
    IF _parent_address IS NOT NULl THEN
        SELECT id INTO stored_parent_id FROM "listed_tokens" WHERE address = _parent_address;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'token with id % does not exist', stored_parent_id;
        END IF;
    END IF;
    
    SELECT * INTO
        stored_token
    FROM "listed_tokens"
    WHERE address = _address;
    IF NOT FOUND THEN
        INSERT INTO "listed_tokens"(address, name, symbol, timestamp, parent_id)
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
                UPDATE "listed_tokens"
                SET name=_name,
                    symbol=_symbol,
                    timestamp=_timestamp,
                    parent_id = stored_parent_id
                WHERE address = _address RETURNING * INTO stored_token;
                changed := TRUE;
            END IF;
        END IF;
    END IF;


    INSERT INTO "listed_tokens_reserves"(address)
    VALUES (_reserve_address)
    ON CONFLICT DO NOTHING;
    SELECT * INTO stored_reserve FROM "listed_tokens_reserves" WHERE address = _reserve_address;

    INSERT INTO "listed_tokens_reserves_tokens"(token_id, reserve_id)
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
      FROM "listed_tokens" AS toks
               LEFT JOIN "listed_tokens" AS olds
                         ON toks.id = olds.parent_id
               JOIN "%[4]s" as token_reserve 
                         ON toks.id = token_reserve.token_id
               JOIN "%[3]s" as reserve
                         ON reserve.id = token_reserve.reserve_id
      WHERE toks.parent_id IS NULL
      ORDER BY timestamp DESC) AS joined
GROUP BY joined.address, joined.name, joined.symbol, joined.timestamp, joined.reserve_address;
`
