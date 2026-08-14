-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian2exchange (

    idexchange BIGINT NOT NULL,
    idcustodian BIGINT NOT NULL,

    value01 VARCHAR(255) NULL,
    value02 SMALLINT DEFAULT 0 NOT NULL,

    PRIMARY KEY (idexchange, idcustodian),

    CONSTRAINT fk_exchange FOREIGN KEY (idexchange)
        REFERENCES exchange(idexchange) ON DELETE CASCADE,

    CONSTRAINT fk_custodian FOREIGN KEY (idcustodian)
        REFERENCES custodian(idcustodian) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian2exchange
-- +goose StatementEnd

