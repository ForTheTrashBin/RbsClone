-- +goose up
-- +goose StatementBegin
CREATE TABLE exchange (

    idexchange UUID PRIMARY KEY DEFAULT uuidv7(),

    shortcode VARCHAR(8) NOT NULL UNIQUE,

    lastname VARCHAR(80) NOT NULL,
    firstname VARCHAR(80) NULL,

    statuscode SMALLINT DEFAULT 0 NOT NULL,
    scorepoints SMALLINT DEFAULT 0 NOT NULL,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$')
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE exchange
-- +goose StatementEnd
