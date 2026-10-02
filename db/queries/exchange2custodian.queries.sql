-- ============================================================================
-- Table: exchange2custodian
-- ============================================================================

-- name: GetExchange2CustodianByIdcustodianAndIdexchange :one

SELECT * FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1 AND idexchange = $2 LIMIT 1;

-- name: GetExchange2CustodianByIdcustodian :many

SELECT * FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1 ORDER BY sequenceno;

-- name: GetExchange2CustodianCountByIdcustodian :one

SELECT COUNT(*) FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1;

-- ----------------------------------------------------------------------------

-- name: InsertExchange2Custodian :exec

INSERT INTO EXCHANGE2CUSTODIAN (idcustodian, idexchange, sequenceno, value1, value2) VALUES ($1, $2, $3, $4, $5);

-- ----------------------------------------------------------------------------

-- name: DeleteExchange2CustodianByIdcustodianAndIdexchange :execresult

DELETE FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1 AND idexchange = $2;

-- name: DeleteExchange2CustodianByIdcustodian :execresult

DELETE FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1;
