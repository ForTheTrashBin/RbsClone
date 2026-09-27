-- +goose up
-- +goose StatementBegin
CREATE TABLE country (

    id UUID PRIMARY KEY DEFAULT uuidv7(),

    shortcode VARCHAR(2) NOT NULL UNIQUE,

    name VARCHAR(80) NOT NULL,
    
    flags SMALLINT DEFAULT 0 NOT NULL,
    
    ibanlength SMALLINT DEFAULT 22 NULL,

    risktype SMALLINT DEFAULT 0 NOT NULL,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$')
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE country
-- +goose StatementEnd

