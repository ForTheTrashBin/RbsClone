-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian2exchange (

    exchange_id INTEGER NOT NULL,
    custodian_id INTEGER NOT NULL,

    descrption TEXT NULL,
    hours_allocated INTEGER DEFAULT 0 NOT NULL,

    PRIMARY KEY (exchange_id, custodian_id),

    FOREIGN KEY (exchange_id) REFERENCES exchange(id) ON DELETE CASCADE,
    FOREIGN KEY (custodian_id) REFERENCES custodian(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian2exchange
-- +goose StatementEnd
