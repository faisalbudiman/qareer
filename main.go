package main

import (
	"fmt"
	"net/http"
	"os"

	"qareer/internal/locations"
	"qareer/pkg/db"
	internalMiddleware "qareer/pkg/middleware"
	"qareer/pkg/utils"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/facebookgo/grace/gracehttp"
)

var (
	PORT        = utils.GetEnv("PORT", "8000")
	DB_DRIVER   = "pgx"
	DB_HOST     = utils.GetEnv("DB_HOST", "db")
	DB_USER     = utils.GetEnv("DB_USER", "qareer")
	DB_PASSWORD = utils.GetEnv("DB_PASSWORD", "qareer")
	DB_NAME     = utils.GetEnv("DB_NAME", "qareer")
	logger      *log.Logger
	connector   db.Db
)

func init() {
	logger = log.New()
	// logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)
	connector = connector.Constructor(DB_DRIVER, fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		DB_USER, DB_PASSWORD, DB_HOST, DB_NAME))
}

func main() {
	r := chi.NewRouter()
	// middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(internalMiddleware.LoggerMiddleware(&middleware.DefaultLogFormatter{
		Logger: logger.WithField("type", "access_log"),
	}, r))

	// router
	r.Mount("/locations", locations.DefaultService(locations.ServiceConfig{
		Db:       connector,
		Logger:   logger,
		Response: new(utils.HTTPJSONResponse),
	}))

	println("serve at", PORT)
	gracehttp.Serve(
		&http.Server{Addr: "0.0.0.0" + PORT, Handler: r},
	)
}
