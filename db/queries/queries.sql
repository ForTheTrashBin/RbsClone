-- ============================================================================
-- Table: custodian
-- ============================================================================

-- name: CreateCustodian :one

INSERT INTO CUSTODIAN (shortcode, lastname, firstname, statuscode, scorepoints) VALUES ($1, $2, $3, $4, $5) RETURNING idCustodian;

-- name: CreateCustodianMin :one

INSERT INTO CUSTODIAN (shortcode, lastname) VALUES ($1, $2) RETURNING idCustodian;

-- ----------------------------------------------------------------------------

-- name: ListCustodians :many

SELECT * FROM CUSTODIAN ORDER BY shortcode;

-- name: GetCustodianByID :one

SELECT * FROM CUSTODIAN WHERE idCustodian = $1 LIMIT 1;

-- name: GetCustodianByShortcode :one

SELECT * FROM CUSTODIAN WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian :exec

UPDATE CUSTODIAN SET shortcode = $2, lastname = $3, firstname = $4, statuscode = $5, scorepoints = $6 WHERE idcustodian = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteCustodianByIdCustodian :exec

DELETE FROM CUSTODIAN WHERE idcustodian = $1;

-- name: DeleteCustodianByShortcode :exec

DELETE FROM CUSTODIAN WHERE shortcode = $1;

-- ============================================================================
-- Table: exchange
-- ============================================================================

-- name: CreateExchange :one

INSERT INTO EXCHANGE (shortcode, lastname, firstname, statuscode, scorepoints) VALUES ($1, $2, $3, $4, $5) RETURNING idExchange;

-- name: CreateExchangeMin :one

INSERT INTO EXCHANGE (shortcode, lastname) VALUES ($1, $2) RETURNING idExchange;

-- ----------------------------------------------------------------------------

-- name: ListExchanges :many

SELECT * FROM EXCHANGE ORDER BY shortcode;

-- name: GetExchangeByID :one

SELECT * FROM EXCHANGE WHERE idExchange = $1 LIMIT 1;

-- name: GetExchangeByShortcode :one

SELECT * FROM EXCHANGE WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateExchange :exec

UPDATE EXCHANGE SET shortcode = $2, lastname = $3, firstname = $4, statuscode = $5, scorepoints = $6 WHERE idexchange = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteExchangeByIdExchange :exec

DELETE FROM EXCHANGE WHERE idexchange = $1;

-- name: DeleteExchangeByShortcode :exec

DELETE FROM EXCHANGE WHERE shortcode = $1;

-- ============================================================================
-- Table: custodian2exchange
-- ============================================================================

-- name: CreateCustodian2Exchange :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, value01, value02) VALUES ($1, $2, $3, $4);

-- name: CreateCustodian2ExchangeMin :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, value02) VALUES ($1, $2, $3);

-- ----------------------------------------------------------------------------

-- name: ListCustodian2ExchangeByExchange :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 ORDER BY idcustodian;

-- name: ListCustodian2ExchangeByCustodian :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idcustodian = $1 ORDER BY idexchange;

-- name: GetCustodian2ExchangeByExchangeAndCustodian :one

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 AND idcustodian = $2 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian2Exchange :exec

UPDATE CUSTODIAN2EXCHANGE SET value01 = $2, value02 = $3 WHERE idexchange = $1 AND idcustodian = $2;

-- ----------------------------------------------------------------------------

-- name: DeleteExchangeByIdExchangeAndIdCustodian :exec

DELETE FROM CUSTODIAN2EXCHANGE WHERE idexchange = $1 AND idcustodian = $2;
