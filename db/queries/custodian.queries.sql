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
