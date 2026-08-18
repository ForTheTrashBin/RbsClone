-- ============================================================================
-- Table: custodian2exchange
-- ============================================================================

-- name: InsertCustodian2Exchange :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, flags, value01, value02) VALUES ($1, $2, $3, $4, $5);

-- ----------------------------------------------------------------------------

-- name: GetCustodian2ExchangeByIdExchange :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 ORDER BY idcustodian;

-- name: GetCustodian2ExchangeByIdCustodian :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idcustodian = $1 ORDER BY idexchange;

-- name: GetCustodian2ExchangeByIdExchangeAndIdCustodian :one

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 AND idcustodian = $2 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian2Exchange :execresult

UPDATE CUSTODIAN2EXCHANGE SET flags = $3, value01 = $4, value02 = $5,  WHERE idexchange = $1 AND idcustodian = $2;

-- ----------------------------------------------------------------------------

-- name: DeleteCustodian2Exchange :execresult

DELETE FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 AND idcustodian = $2;
