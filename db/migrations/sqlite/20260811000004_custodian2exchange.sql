-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian2exchange (

    idexchange INTEGER NOT NULL,
    idcustodian INTEGER NOT NULL,

    value01 TEXT NULL,
    value02 INTEGER DEFAULT 0 NOT NULL,

    PRIMARY KEY (idexchange, idcustodian),

    FOREIGN KEY (idexchange) REFERENCES exchange(idexchange) ON DELETE CASCADE,
    FOREIGN KEY (idcustodian) REFERENCES custodian(idcustodian) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian2exchange
-- +goose StatementEnd
