-- +goose up
-- +goose StatementBegin
CREATE TABLE country (

    idcountry BIGINT AUTO_INCREMENT PRIMARY KEY,

    shortcode VARCHAR(2) NOT NULL UNIQUE,

    CONSTRAINT chk_shortcode_uppercase CHECK (REGEXP_LIKE(shortcode, '^[A-Z0-9ÄÖÜß]+$', 'c')),

    name VARCHAR(30) NOT NULL,

    flags SMALLINT DEFAULT 0 NULL,

    ibanlength SMALLINT DEFAULT 22 NULL,

    risktype SMALLINT DEFAULT 0 NULL
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd
