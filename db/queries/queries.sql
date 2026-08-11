-- name: CreateExchange :one

INSERT INTO exchange (user_code, last_name, first_name, status_code, score_points)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_code, last_name, first_name, status_code, score_points;

-- name: GetExchange :one

SELECT * FROM EXCHANGE
WHERE user_code = $1 LIMIT 1;

-- name: ListExchanges :many

SELECT * FROM EXCHANGE
ORDER BY user_code;

