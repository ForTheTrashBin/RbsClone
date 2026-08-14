-- +goose up
-- +goose StatementBegin
CREATE TABLE country (

    idcountry BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    shortcode VARCHAR(2) NOT NULL UNIQUE,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$'),

    name VARCHAR(30) NOT NULL,
    
    flags SMALLINT DEFAULT 0 NOT NULL,
    
    ibanlength SMALLINT DEFAULT 22 NULL,

    risktype SMALLINT DEFAULT 0 NOT NULL
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE country
-- +goose StatementEnd

