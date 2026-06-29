package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

var (
	echoInstance *echo.Echo
	initOnce     sync.Once
	initErr      error
)

func initApp() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		initErr = log.Output(2, "DATABASE_URL is not set")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		initErr = err
		return
	}

	if err := db.Ping(); err != nil {
		initErr = err
		return
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/openapi.json", func(c echo.Context) error {
		spec, err := generated.GetSwagger()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, spec)
	})

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

	swagger, err := generated.GetSwagger()
	if err != nil {
		initErr = err
		return
	}
	e.Use(oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			return path == "/docs" || path == "/openapi.json"
		},
		SilenceServersWarning: true,
	}))

	repo := repository.NewPostgresRepository(db)
	server := handler.NewServer(repo)
	generated.RegisterHandlers(e, server)

	echoInstance = e
}

// Handler is the entry point for Vercel serverless functions.
func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(initApp)
	if initErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	echoInstance.ServeHTTP(w, r)
}
