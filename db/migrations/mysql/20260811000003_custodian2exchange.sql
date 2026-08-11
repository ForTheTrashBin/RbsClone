-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian2exchange (

    exchange_id BIGINT NOT NULL,
    custodian_id BIGINT NOT NULL,

    descrption VARCHAR(255) NULL,
    hours_allocated SMALLINT DEFAULT 0 NOT NULL,

    PRIMARY KEY (exchange_id, custodian_id),

    CONSTRAINT fk_exchange FOREIGN KEY (exchange_id)
        REFERENCES exchange(id) ON DELETE CASCADE,

    CONSTRAINT fk_custodian FOREIGN KEY (custodian_id)
        REFERENCES custodian(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian2exchange
-- +goose StatementEnd
