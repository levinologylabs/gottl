// Package db provides the database queries and utiltiies around migrations and transactions. It is
// in part, genrated code using sqlc based off the *.sql files in the same directory.
//
// This package also provides the QueriesExt struct which extends the functionality of the generated
// Queries struct to include the ability to work with transactions and the underlying pgx connection.
package db
