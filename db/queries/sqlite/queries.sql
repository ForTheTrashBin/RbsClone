-- ============================================================================
-- Table: country
-- ============================================================================

-- name: InsertCountry :one

INSERT INTO COUNTRY (shortcode, name, flags, ibanlength, risktype) VALUES (?, ?, ?, ?, ?) RETURNING idCountry;

-- name: InsertCountryMin :one

INSERT INTO COUNTRY (shortcode, name) VALUES (?, ?) RETURNING idCountry;

-- ----------------------------------------------------------------------------

-- name: ListCountries :many

SELECT * FROM COUNTRY ORDER BY shortcode;

-- name: GetCountryByID :one

SELECT * FROM COUNTRY WHERE idCountry = ? LIMIT 1;

-- name: GetCountryByShortcode :one

SELECT * FROM COUNTRY WHERE shortcode = ? LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCountry :exec

UPDATE COUNTRY SET shortcode = ?, name = ?, flags = ?, ibanlength = ?, risktype = ? WHERE idcountry = ?;

-- ----------------------------------------------------------------------------

-- name: DeleteCountryByIdCountry :exec

DELETE FROM COUNTRY WHERE idcountry = ?;

-- name: DeleteCountryByShortcode :exec

DELETE FROM COUNTRY WHERE shortcode = ?;

-- ============================================================================
-- Table: custodian
-- ============================================================================

-- name: InsertCustodian :one

INSERT INTO CUSTODIAN (shortcode, name, flags, idcountry, depotno) VALUES (?, ?, ?, ?, ?) RETURNING idCustodian;

-- name: InsertCustodianMin :one

INSERT INTO CUSTODIAN (shortcode, name, idcountry) VALUES (?, ?, ?) RETURNING idCustodian;

-- ----------------------------------------------------------------------------

-- name: ListCustodians :many

SELECT * FROM CUSTODIAN ORDER BY shortcode;

-- name: GetCustodianByID :one

SELECT * FROM CUSTODIAN WHERE idCustodian = ? LIMIT 1;

-- name: GetCustodianByShortcode :one

SELECT * FROM CUSTODIAN WHERE shortcode = ? LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian :exec

UPDATE CUSTODIAN SET shortcode = ?, name = ?, flags = ?, idcountry = ?, depotno = ? WHERE idcustodian = ?;

-- ----------------------------------------------------------------------------

-- name: DeleteCustodianByIdCustodian :exec

DELETE FROM CUSTODIAN WHERE idcustodian = ?;

-- name: DeleteCustodianByShortcode :exec

DELETE FROM CUSTODIAN WHERE shortcode = ?;

-- ============================================================================
-- Table: exchange
-- ============================================================================

-- name: InsertExchange :one

INSERT INTO EXCHANGE (shortcode, lastname, firstname, statuscode, scorepoints) VALUES (?, ?, ?, ?, ?) RETURNING idExchange;

-- name: InsertExchangeMin :one

INSERT INTO EXCHANGE (shortcode, lastname) VALUES (?, ?) RETURNING idExchange;

-- ----------------------------------------------------------------------------

-- name: ListExchanges :many

SELECT * FROM EXCHANGE ORDER BY shortcode;

-- name: GetExchangeByID :one

SELECT * FROM EXCHANGE WHERE idExchange = ? LIMIT 1;

-- name: GetExchangeByShortcode :one

SELECT * FROM EXCHANGE WHERE shortcode = ? LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateExchange :exec

UPDATE EXCHANGE SET shortcode = ?, lastname = ?, firstname = ?, statuscode = ?, scorepoints = ? WHERE idexchange = ?;

-- ----------------------------------------------------------------------------

-- name: DeleteExchangeByIdExchange :exec

DELETE FROM EXCHANGE WHERE idexchange = ?;

-- name: DeleteExchangeByShortcode :exec

DELETE FROM EXCHANGE WHERE shortcode = ?;

-- ============================================================================
-- Table: custodian2exchange
-- ============================================================================

-- name: InsertCustodian2Exchange :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, value01, value02) VALUES (?, ?, ?, ?);

-- name: InsertCustodian2ExchangeMin :exec

INSERT INTO CUSTODIAN2EXCHANGE (idexchange, idcustodian, value02) VALUES (?, ?, ?);

-- ----------------------------------------------------------------------------

-- name: ListCustodian2ExchangeByExchange :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = ? ORDER BY idcustodian;

-- name: ListCustodian2ExchangeByCustodian :many

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idcustodian = ? ORDER BY idexchange;

-- name: GetCustodian2ExchangeByExchangeAndCustodian :one

SELECT * FROM CUSTODIAN2EXCHANGE WHERE idexchange = ? AND idcustodian = ? LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian2Exchange :exec

UPDATE CUSTODIAN2EXCHANGE SET value01 = ?, value02 = ? WHERE idexchange = ? AND idcustodian = ?;

-- ----------------------------------------------------------------------------

-- name: DeleteExchangeByIdExchangeAndIdCustodian :exec

DELETE FROM CUSTODIAN2EXCHANGE WHERE idexchange = ? AND idcustodian = ?;
