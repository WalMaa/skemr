package main

import (
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"skemr/db/sqlc"
	"skemr/routers"
	"skemr/service"
)

// runSchema drops the current schema, reads schema.sql file and executes it to set up the database schema.
func runSchema(conn *pgx.Conn) {
	schema, err := ioutil.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Exec(context.Background(), "DROP SCHEMA IF EXISTS public CASCADE")
	_, err = conn.Exec(context.Background(), string(schema))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:pass@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(conn)
	projectService := service.NewProjectService(queries)
	databaseService := service.NewDatabaseService(queries)

	runSchema(conn)

	// Initialize services
	services := &routers.Services{
		ProjectService:  projectService,
		DatabaseService: databaseService,
	}

	r := routers.InitRouter(services)
	defer conn.Close(context.Background())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
