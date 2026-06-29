package main

import (
	"database/sql"
	"log"
	"net/http"
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

	// Serve OpenAPI spec as JSON (registered before validator so it's not blocked)
	e.GET("/openapi.json", func(c echo.Context) error {
		spec, err := generated.GetSwagger()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, spec)
	})

	// Serve Swagger UI at /docs (registered before validator so it's not blocked)
	e.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<!DOCTYPE html>
<html>
<head>
    <title>Estate Service - API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/openapi.json",
            dom_id: '#swagger-ui',
            presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
            layout: "BaseLayout"
        })
    </script>
</body>
</html>`)
	})

	// OpenAPI validator (skips /docs and /openapi.json)
	swagger, err := generated.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	e.Use(oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			return path == "/docs" || path == "/openapi.json"
		},
		SilenceServersWarning: true,
	}))

	// Repository
	repo := repository.NewPostgresRepository(db)

	// Server
	server := handler.NewServer(repo)

	// Register API routes
	generated.RegisterHandlers(e, server)

	log.Println("Server started on :1323")

	e.Logger.Fatal(
		e.Start(":1323"),
	)
}
