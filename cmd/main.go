package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // Ensure the database connection is closed when the application exits

	if err := db.Ping(); err != nil {
		log.Fatal(err) // quit if the database connection cannot be established
	}

	// Echo
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// OpenAPI validator
	swagger, err := generated.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(oapimiddleware.OapiRequestValidator(swagger))

	// Repository
	repo := repository.NewPostgresRepository(db)

	// Server
	server := handler.NewServer(repo)

	// Register routes
	generated.RegisterHandlers(e, server)

	//e.GET("/swagger/*", swagger.WrapHandlerV3)

	log.Println("Server started on :1323")

	e.Logger.Fatal(
		e.Start(":1323"),
	)
}
