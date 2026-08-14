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
