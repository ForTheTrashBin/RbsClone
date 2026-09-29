-- ============================================================================
-- Table: custodian2exchange
-- ============================================================================

-- name: GetCustodian2ExchangeByIdexchangeAndIdcustodian :one

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idcustodian = $1 AND idexchange = $2 LIMIT 1;

-- name: GetCustodian2ExchangeByIdcustodian :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idcustodian = $1 ORDER BY sequenceno;

-- ----------------------------------------------------------------------------

-- name: InsertCustodian2Exchange :exec

INSERT INTO CUSTODIAN2EXCHANGE (idcustodian, idexchange, sequenceno) VALUES ($1, $2, $3);

-- ----------------------------------------------------------------------------

-- name: DeleteCustodian2Exchange :execresult

DELETE FROM CUSTODIAN2EXCHANGE WHERE idcustodian = $1 AND idexchange = $2;
