-- +goose up
-- +goose StatementBegin
CREATE TABLE custodian2exchange (

    idcustodian UUID NOT NULL,
    idexchange UUID NOT NULL,
    sequenceno int NOT NULL,

    PRIMARY KEY (idcustodian, idexchange),

    CONSTRAINT cs_idcustodiansequenceno UNIQUE (idcustodian, sequenceno),

    CONSTRAINT fk_exchange  FOREIGN KEY (idexchange)  REFERENCES exchange(id) ON DELETE CASCADE,
    CONSTRAINT fk_custodian FOREIGN KEY (idcustodian) REFERENCES custodian(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose down

-- +goose StatementBegin
DROP TABLE custodian2exchange
-- +goose StatementEnd

