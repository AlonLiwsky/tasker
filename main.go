package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tasker/entities"
	"github.com/tasker/repo"
	"github.com/tasker/service"
	"github.com/tasker/web"
)

func main() {
	//Setup DB
	db, err := setupDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Create repo
	repo := repo.NewRepository(db)

	//Create Step Runners
	apiCallerStepRunner := service.NewApiCallerStepRunner(http.Client{})
	stepRunners := map[entities.StepType]service.StepRunner{
		entities.APICallStepType: apiCallerStepRunner,
	}

	//Create service
	srv := service.NewService(repo, stepRunners)

	//Create adapter
	adapter := web.NewAdapter(srv)

	//Create router
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Accept any origin

	})
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/task", func(r chi.Router) {
		r.Post("/", adapter.CreateTask) // POST /articles
		r.Get("/{taskID}", adapter.GetTask)
		r.Post("/{taskID}/execute/{scheduleID}", adapter.ExecuteTask)
	})

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	http.ListenAndServe(":3333", r)
}

func setupDB() (*sql.DB, error) {
	// Connect to database
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database_name")
	if err != nil {
		return nil, err
	}

	//Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	// Read the SQL file
	file, err := os.Open("tables.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the SQL file
	sqlBytes, err := os.ReadFile("tables.sql")
	if err != nil {
		log.Fatal(err)
	}
	sqlStatements := strings.Split(string(sqlBytes), ";")

	// Execute each SQL statement in the file
	for _, stmt := range sqlStatements {
		trimmedStmt := strings.TrimSpace(stmt)
		if trimmedStmt != "" {
			_, err := db.Exec(trimmedStmt)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

/* TO DO:
-Continue with get task endpoint
-Handle idempotency
-...
*/
