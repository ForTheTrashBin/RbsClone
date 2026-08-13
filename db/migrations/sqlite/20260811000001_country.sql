-- +goose up
-- +goose StatementBegin
CREATE TABLE country (

    idcountry INTEGER PRIMARY KEY AUTOINCREMENT,

    shortcode TEXT NOT NULL UNIQUE,

    CONSTRAINT chk_shortcode_rules CHECK (
        length(shortcode) <= 8 AND
        shortcode GLOB '*[A-Z0-9ÄÖÜß]*' AND
        NOT shortcode GLOB '*[a-zäöü]*'
    ),

    name TEXT NOT NULL,

    flags INTEGER DEFAULT 0,

    ibanlength INTEGER DEFAULT 22,
    
    risktype INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE country
-- +goose StatementEnd
