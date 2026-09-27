-- ============================================================================
-- Table: custodian
-- ============================================================================

-- name: GetCustodians :many

SELECT * FROM CUSTODIAN ORDER BY shortcode;

-- name: GetCustodianById :one

SELECT * FROM CUSTODIAN WHERE id = $1 LIMIT 1;

-- name: GetCustodianByShortcode :one

SELECT * FROM CUSTODIAN WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: InsertCustodian :one

INSERT INTO CUSTODIAN (shortcode, name, flags, idcountry, depotno) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- ----------------------------------------------------------------------------

-- name: UpdateCustodian :one

UPDATE CUSTODIAN SET shortcode = $2, name = $3, flags = $4, idcountry = $5, depotno = $6 WHERE id = $1 RETURNING *;

-- ----------------------------------------------------------------------------

-- name: DeleteCustodian :execresult

DELETE FROM CUSTODIAN WHERE id = $1;
