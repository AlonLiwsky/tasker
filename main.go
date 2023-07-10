package main

import (
	"context"
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
	apicall2 "github.com/tasker/repo/apicall"
	"github.com/tasker/repo/executionDB"
	"github.com/tasker/repo/mgmtDB"
	"github.com/tasker/service"
	"github.com/tasker/service/apicall"
	"github.com/tasker/service/storageread"
	"github.com/tasker/service/storagewrite"
	"github.com/tasker/web"

	"github.com/redis/go-redis/v9"
)

func main() {
	//Setup sql DB
	sqlDB, err := setupMgmtDB()
	if err != nil {
		panic(err.Error())
	}
	defer sqlDB.Close()

	//Setup redis DB
	redis, err := setupExecutionDB()
	if err != nil {
		panic(err.Error())
	}

	//Setup http client
	httpClient := http.Client{}

	//Create repos
	mgmtRepo := mgmtDB.NewRepository(sqlDB)
	executionRepo := executionDB.NewRepository(redis)
	apicallRepo := apicall2.NewRepository(httpClient)

	//Create Step Runners
	apiCallerStepRunner := apicall.NewStepRunner(apicallRepo)
	storageReadStepRunner := storageread.NewStepRunner(executionRepo)
	storageWriteStepRunner := storagewrite.NewStepRunner(executionRepo)
	stepRunners := map[entities.StepType]service.StepRunner{
		entities.APICallStepType:      apiCallerStepRunner,
		entities.StorageReadStepType:  storageReadStepRunner,
		entities.StorageWriteStepType: storageWriteStepRunner,
	}

	//Create service
	srv := service.NewService(mgmtRepo, stepRunners)

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

	r.Route("/schedule", func(r chi.Router) {
		r.Post("/", adapter.CreateSchedule) // POST /articles
	})

	r.Route("/jobs", func(r chi.Router) {
		r.Post("/execute-scheduled-tasks", adapter.ExecuteScheduledTasks) // POST /articles
	})

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	http.ListenAndServe(":3333", r)
}

func setupExecutionDB() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ping the Redis server to check the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func setupMgmtDB() (*sql.DB, error) {
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
