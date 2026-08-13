-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian (

    idcustodian BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    shortcode VARCHAR(5) NOT NULL UNIQUE,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$'),

    name VARCHAR(30) NOT NULL,

    flags SMALLINT DEFAULT 0 NULL,

    idcountry BIGINT NOT NULL,

    CONSTRAINT fk_country FOREIGN KEY (idcountry)
        REFERENCES country(idcountry) ON DELETE RESTRICT,

    depotno VARCHAR(10) NULL
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd
