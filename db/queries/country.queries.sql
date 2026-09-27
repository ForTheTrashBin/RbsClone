-- ============================================================================
-- Table: country
-- ============================================================================

-- name: GetCountries :many

SELECT * FROM COUNTRY ORDER BY shortcode;

-- name: GetCountryById :one

SELECT * FROM COUNTRY WHERE id = $1 LIMIT 1;

-- name: GetCountryByShortcode :one

SELECT * FROM COUNTRY WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: InsertCountry :one

INSERT INTO COUNTRY (shortcode, name, flags, ibanlength, risktype) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- ----------------------------------------------------------------------------

-- name: UpdateCountry :one

UPDATE COUNTRY SET shortcode = $2, name = $3, flags = $4, ibanlength = $5, risktype = $6 WHERE id = $1 RETURNING *;

-- ----------------------------------------------------------------------------

-- name: DeleteCountry :execresult

DELETE FROM COUNTRY WHERE id = $1;
