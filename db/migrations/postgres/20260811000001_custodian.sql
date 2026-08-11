-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian (

    idcustodian BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    shortcode VARCHAR(8) NOT NULL UNIQUE,

    CONSTRAINT chk_shortcode_uppercase CHECK (shortcode ~ '^[A-Z0-9ÄÖÜß]+$'),

    lastname VARCHAR(80) NOT NULL,
    firstname VARCHAR(80) NULL,

    statuscode SMALLINT DEFAULT 0 NOT NULL,
    scorepoints SMALLINT DEFAULT 0 NOT NULL
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd

