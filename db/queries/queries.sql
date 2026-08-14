-- ============================================================================
-- Table: country
-- ============================================================================

-- name: InsertCountry :one

INSERT INTO COUNTRY (shortcode, name, flags, ibanlength, risktype) VALUES ($1, $2, $3, $4, $5) RETURNING idCountry;

-- name: InsertCountryMin :one

INSERT INTO COUNTRY (shortcode, name) VALUES ($1, $2) RETURNING idCountry;

-- ----------------------------------------------------------------------------

-- name: GetCountries :many

SELECT * FROM COUNTRY ORDER BY shortcode;

-- name: GetCountryByID :one

SELECT * FROM COUNTRY WHERE idCountry = $1 LIMIT 1;

-- name: GetCountryByShortcode :one

SELECT * FROM COUNTRY WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCountry :exec

UPDATE COUNTRY SET shortcode = $2, name = $3, flags = $4, ibanlength = $5, risktype = $6 WHERE idcountry = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteCountryById :exec

DELETE FROM COUNTRY WHERE idcountry = $1;

-- name: DeleteCountryByShortcode :exec

DELETE FROM COUNTRY WHERE shortcode = $1;

-- ============================================================================
-- Table: custodian
-- ============================================================================

-- name: InsertCustodian :one

INSERT INTO CUSTODIAN (shortcode, name, flags, idcountry, depotno) VALUES ($1, $2, $3, $4, $5) RETURNING idCustodian;

-- name: InsertCustodianMin :one

INSERT INTO CUSTODIAN (shortcode, name, idcountry) VALUES ($1, $2, $3) RETURNING idCustodian;

-- ----------------------------------------------------------------------------

-- name: GetCustodians :many

SELECT * FROM CUSTODIAN ORDER BY shortcode;

-- name: GetCustodianByID :one

SELECT * FROM CUSTODIAN WHERE idCustodian = $1 LIMIT 1;

-- name: GetCustodianByShortcode :one

SELECT * FROM CUSTODIAN WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian :exec

UPDATE CUSTODIAN SET shortcode = $2, name = $3, flags = $4, idcountry = $5, depotno = $6 WHERE idcustodian = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteCustodianById :exec

DELETE FROM CUSTODIAN WHERE idcustodian = $1;

-- name: DeleteCustodianByShortcode :exec

DELETE FROM CUSTODIAN WHERE shortcode = $1;

-- ============================================================================
-- Table: exchange
-- ============================================================================

-- name: InsertExchange :one

INSERT INTO EXCHANGE (shortcode, lastname, firstname, statuscode, scorepoints) VALUES ($1, $2, $3, $4, $5) RETURNING idExchange;

-- name: InsertExchangeMin :one

INSERT INTO EXCHANGE (shortcode, lastname) VALUES ($1, $2) RETURNING idExchange;

-- ----------------------------------------------------------------------------

-- name: GetExchanges :many

SELECT * FROM EXCHANGE ORDER BY shortcode;

-- name: GetExchangeByID :one

SELECT * FROM EXCHANGE WHERE idExchange = $1 LIMIT 1;

-- name: GetExchangeByShortcode :one

SELECT * FROM EXCHANGE WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateExchange :exec

UPDATE EXCHANGE SET shortcode = $2, lastname = $3, firstname = $4, statuscode = $5, scorepoints = $6 WHERE idexchange = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteExchangeById :exec

DELETE FROM EXCHANGE WHERE idexchange = $1;

-- name: DeleteExchangeByShortcode :exec

DELETE FROM EXCHANGE WHERE shortcode = $1;

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
