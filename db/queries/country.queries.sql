-- ============================================================================
-- Table: country
-- ============================================================================

-- name: InsertCountry :one

INSERT INTO COUNTRY (shortcode, name, flags, ibanlength, risktype) VALUES ($1, $2, $3, $4, $5) RETURNING idCountry;

-- ----------------------------------------------------------------------------

-- name: GetCountries :many

SELECT * FROM COUNTRY ORDER BY shortcode;

-- name: GetCountryByID :one

SELECT * FROM COUNTRY WHERE idCountry = $1 LIMIT 1;

-- name: GetCountryByShortcode :one

SELECT * FROM COUNTRY WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateCountry :execresult

UPDATE COUNTRY SET shortcode = $2, name = $3, flags = $4, ibanlength = $5, risktype = $6 WHERE idcountry = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteCountry :execresult

DELETE FROM COUNTRY WHERE idcountry = $1;
