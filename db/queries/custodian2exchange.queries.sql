-- ============================================================================
-- Table: custodian2exchange
-- ============================================================================

-- name: InsertCustodian2Exchange :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, value01, value02) VALUES ($1, $2, $3, $4);

-- name: InsertCustodian2ExchangeMin :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, value02) VALUES ($1, $2, $3);

-- ----------------------------------------------------------------------------

-- name: GetCustodian2ExchangeByExchange :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 ORDER BY idcustodian;

-- name: GetCustodian2ExchangeByCustodian :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idcustodian = $1 ORDER BY idexchange;

-- name: GetCustodian2ExchangeByExchangeAndCustodian :one

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 AND idcustodian = $2 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian2Exchange :exec

UPDATE CUSTODIAN2EXCHANGE SET value01 = $2, value02 = $3 WHERE idexchange = $1 AND idcustodian = $2;

-- ----------------------------------------------------------------------------

-- name: DeleteCustodian2Exchange :exec

DELETE FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 AND idcustodian = $2;
