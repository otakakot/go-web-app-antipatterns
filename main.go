package main

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/otakakot/go-web-app-antipatterns/internal/app"
	"github.com/otakakot/go-web-app-antipatterns/internal/handler"
	"github.com/otakakot/go-web-app-antipatterns/internal/repository"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err := repository.MigrateDB(db); err != nil {
		panic(err)
	}

	txProvider := repository.NewTransactionProvider(db)

	usePointsAsDiscountHandler := app.NewUsePointsAsDiscountHandler(txProvider)

	handler := handler.NewHTTPHandler(usePointsAsDiscountHandler)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
