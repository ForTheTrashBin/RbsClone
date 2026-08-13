-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian (

    idcustodian INTEGER PRIMARY KEY AUTOINCREMENT,

    shortcode TEXT NOT NULL UNIQUE,

    CONSTRAINT chk_shortcode_rules CHECK (
        length(shortcode) <= 8 AND
        shortcode GLOB '*[A-Z0-9ÄÖÜß]*' AND
        NOT shortcode GLOB '*[a-zäöü]*'
    ),

    name TEXT NOT NULL,

    flags INTEGER DEFAULT 0,

    idcountry INTEGER NOT NULL,

    FOREIGN KEY (idcountry) REFERENCES country(idcountry) ON DELETE CASCADE,

    depotno TEXT
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd
