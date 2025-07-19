package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"skemr/controller"
	"skemr/db/sqlc"
	"skemr/service"
)

// runSchema drops the current schema, reads schema.sql file and executes it to set up the database schema.
func runSchema(conn *pgx.Conn) {
	schema, err := ioutil.ReadFile("schema.sql")
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
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Use(gin.BasicAuth(gin.Accounts{
		"user": "pass",
	}))

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:pass@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(conn)
	projectService := service.NewProjectService(queries)
	projectHandler := controller.NewProjectController(projectService)
	databaseService := service.NewDatabaseService(queries)
	databaseHandler := controller.NewDatabaseController(databaseService)
	databaseHandler.RegisterRoutes(r)
	projectHandler.RegisterRoutes(r)

	runSchema(conn)

	defer conn.Close(context.Background())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
