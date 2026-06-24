-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
)
VALUES (
     $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
FOR NO KEY UPDATE; -- Lock Database but still allow other transaction to query the data

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = $1
RETURNING *;
-- In commonly UPDATE will lock the database but still allow other transaction to query the data
-- So we can use it directly instead of use SELECT ... FOR UPDATE and then UPDATE

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;