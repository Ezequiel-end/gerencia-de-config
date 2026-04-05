package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"api-go/user"

	_ "github.com/lib/pq"
)

func main() {
	repo, cleanup := buildRepository()
	if cleanup != nil {
		defer cleanup()
	}

	service := user.NewUserService(repo)
	controller := user.NewUserController(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListUsers(w, r)
		case http.MethodPost:
			controller.CreateUser(w, r)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUser(w, r)
		case http.MethodPut:
			controller.UpdateUser(w, r)
		case http.MethodDelete:
			controller.DeleteUser(w, r)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func buildRepository() (user.UserRepository, func()) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set; using in-memory repository")
		return user.NewUserRepository(), nil
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := user.MigratePostgres(db); err != nil {
		log.Fatal(err)
	}
	log.Println("Using PostgreSQL repository")
	return user.NewPostgresUserRepository(db), func() { _ = db.Close() }
}
