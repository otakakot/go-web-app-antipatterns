package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err := MigrateDB(db); err != nil {
		panic(err)
	}

	txProvider := NewTransactionProvider(db)

	usePointsAsDiscountHandler := NewUsePointsAsDiscountHandler(txProvider)

	handler := NewHTTPHandler(usePointsAsDiscountHandler)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
