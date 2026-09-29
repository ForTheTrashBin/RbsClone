-- ============================================================================
-- Table: exchange2custodian
-- ============================================================================

-- name: GetExchange2CustodianByIdcustodianAndIdexchange :one

SELECT * FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1 AND idexchange = $2 LIMIT 1;

-- name: GetExchange2CustodianByIdcustodian :many

SELECT * FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1 ORDER BY sequenceno;

-- ----------------------------------------------------------------------------

-- name: InsertExchange2Custodian :exec

INSERT INTO EXCHANGE2CUSTODIAN (idcustodian, idexchange, sequenceno) VALUES ($1, $2, $3);

-- ----------------------------------------------------------------------------

-- name: DeleteExchange2Custodian :execresult

DELETE FROM EXCHANGE2CUSTODIAN WHERE idcustodian = $1 AND idexchange = $2;
