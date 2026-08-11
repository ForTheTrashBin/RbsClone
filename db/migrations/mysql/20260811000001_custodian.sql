-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian (

    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    user_code VARCHAR(8) NOT NULL UNIQUE,

    CONSTRAINT chk_user_code_uppercase CHECK (REGEXP_LIKE(user_code ~ '^[A-Z0-9ÄÖÜß]+$', 'c')),

    last_name VARCHAR(80) NOT NULL,
    first_name VARCHAR(80) NULL,

    status_code SMALLINT DEFAULT 0 NOT NULL,
    score_points SMALLINT DEFAULT 0 NOT NULL
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian
-- +goose StatementEnd

