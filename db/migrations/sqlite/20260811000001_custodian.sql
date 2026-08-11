-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian (

    id INTEGER PRIMARY KEY AUTOINCREMENT,

    user_code TEXT NOT NULL UNIQUE,

    last_name TEXT NOT NULL,
    first_name TEXT NULL,

    status_code INTEGER DEFAULT 0 NOT NULL,
    score_points INTEGER DEFAULT 0 NOT NULL,

    CONSTRAINT chk_user_code_rules CHECK (
        length(user_code) <= 8 AND
        user_code GLOB '*[A-Z0-9ÄÖÜß]*' AND
        NOT user_code GLOB '*[a-zäöü]*'
    )
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd
