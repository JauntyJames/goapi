package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gorilla/mux"
)

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Springf("user=%s password=%s dbname=%s"-user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter
}

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	a.Run(":8080")
}
