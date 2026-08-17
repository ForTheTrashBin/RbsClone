-- +goose up
-- +goose StatementBegin
CREATE TABLE exchange (

    idexchange UUID PRIMARY KEY DEFAULT uuidv7(),

    shortcode VARCHAR(8) NOT NULL UNIQUE,

    name VARCHAR(80) NOT NULL,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$')
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE exchange
-- +goose StatementEnd
