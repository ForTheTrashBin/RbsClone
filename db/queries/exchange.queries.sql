-- ============================================================================
-- Table: exchange
-- ============================================================================

-- name: GetExchanges :many

SELECT * FROM EXCHANGE ORDER BY shortcode;

-- name: GetExchangeById :one

SELECT * FROM EXCHANGE WHERE id = $1 LIMIT 1;

-- name: GetExchangeByShortcode :one

SELECT * FROM EXCHANGE WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: InsertExchange :one

INSERT INTO EXCHANGE (shortcode, name, flags) VALUES ($1, $2, $3) RETURNING *;

-- ----------------------------------------------------------------------------

-- name: UpdateExchange :one

UPDATE EXCHANGE SET shortcode = $2, name = $3, flags = $4 WHERE id = $1 RETURNING *;

-- ----------------------------------------------------------------------------

-- name: DeleteExchange :execresult

DELETE FROM EXCHANGE WHERE id = $1;
