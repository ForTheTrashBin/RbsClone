-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian (

    idcustodian UUID PRIMARY KEY DEFAULT uuidv7(),

    shortcode VARCHAR(5) NOT NULL UNIQUE,

    name VARCHAR(30) NOT NULL,

    flags SMALLINT DEFAULT 0 NULL,

    idcountry UUID NOT NULL,

    depotno VARCHAR(10) NULL,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$'),

    CONSTRAINT fk_country FOREIGN KEY (idcountry) REFERENCES country(idcountry) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd
