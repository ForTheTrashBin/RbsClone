-- ============================================================================
-- Table: exchange
-- ============================================================================

-- name: InsertExchange :one

INSERT INTO EXCHANGE (shortcode, name, flags) VALUES ($1, $2, $3) RETURNING idExchange;

-- ----------------------------------------------------------------------------

-- name: GetExchanges :many

SELECT * FROM EXCHANGE ORDER BY shortcode;

-- name: GetExchangeByID :one

SELECT * FROM EXCHANGE WHERE idExchange = $1 LIMIT 1;

-- name: GetExchangeByShortcode :one

SELECT * FROM EXCHANGE WHERE shortcode = $1 LIMIT 1;

-- ----------------------------------------------------------------------------

-- name: UpdateExchange :execresult

UPDATE EXCHANGE SET shortcode = $2, name = $3, flags = $4 WHERE idexchange = $1;

-- ----------------------------------------------------------------------------

-- name: DeleteExchange :execresult

DELETE FROM EXCHANGE WHERE idexchange = $1;
