-- +goose up
-- +goose StatementBegin
CREATE TABLE exchange (

    idexchange INTEGER PRIMARY KEY AUTOINCREMENT,

    shortcode TEXT NOT NULL UNIQUE,

    lastname TEXT NOT NULL,
    firstname TEXT NULL,

    statuscode INTEGER DEFAULT 0 NOT NULL,
    scorepoints INTEGER DEFAULT 0 NOT NULL,

    CONSTRAINT chk_shortcode_rules CHECK (
        length(shortcode) <= 8 AND
        shortcode GLOB '*[A-Z0-9ÄÖÜß]*' AND
        NOT shortcode GLOB '*[a-zäöü]*'
    )
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE exchange
-- +goose StatementEnd
