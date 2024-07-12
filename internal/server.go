package api

import (
	"fmt"
	"net/http"
	"time"
	_ "timeTracker/api/docs"
	"timeTracker/config"
	"timeTracker/internal/handlers"
	"timeTracker/internal/middlewares"
	"timeTracker/internal/storage"
	"timeTracker/pkg/httputils"
	"timeTracker/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
)

type AppServer struct {
	Port        string
	Host        string
	Env         string
	AuthService string
	handlers.Handlers
}

func (app *AppServer) Run(appConfig config.ApiEnvConfig) {
	app.Port = appConfig.Port
	app.Host = appConfig.Host
	app.Env = appConfig.Env
	app.AuthService = appConfig.AuthService
	app.EnvBox = appConfig
	app.Sender = &httputils.Sender{
		Render: render.New(render.Options{
			IndentJSON: true,
		}),
	}

	storage, err := storage.NewPostgresDB()
	if err != nil {
		logger.Log.Error(err)
		panic(err.Error())
	}

	if err := storage.MigratePostgres(); err != nil {
		logger.Log.Fatal(err)
	}

	app.Storage = storage

	router := mux.NewRouter().StrictSlash(true)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.NotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(app.NotFoundHandler)
	router.Methods("GET").Path("/api/users").HandlerFunc(app.GetUsers)
	router.Methods("POST").Path("/api/users").HandlerFunc(app.AddUser)
	router.Methods("PUT").Path("/api/users").HandlerFunc(app.UpdateUser)
	router.Methods("DELETE").Path("/api/users").HandlerFunc(app.DeleteUser)
	router.Methods("POST").Path("/api/start-track").HandlerFunc(app.StartTrackTask)
	router.Methods("POST").Path("/api/end-track").HandlerFunc(app.EndTrackTask)
	router.Methods("POST").Path("/api/calc-time").HandlerFunc(app.CalcTime)

	if app.Env != config.PROD_ENV {
		router.Methods("GET").PathPrefix("/api/docs/").Handler(httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprint("http://localhost:", app.Port, "/api/docs/doc.json")),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		))
	}

	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      app.Env == "DEV",
		ContentTypeNosniff: true,
		SSLRedirect:        true,
	})

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(negroni.HandlerFunc(middlewares.TrackRequestMiddleware))
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allows all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	wrappedRouter := corsMiddleware.Handler(router)
	n.UseHandler(wrappedRouter)

	startupMessage := "Starting API server " + app.Host + " "
	startupMessage = startupMessage + " on port " + app.Port
	startupMessage = startupMessage + " in " + app.Env + " mode."
	logger.Log.Info(startupMessage)

	addr := ":" + app.Port
	if app.Env == "DEV" {
		addr = app.Host + ": " + app.Port
	}

	server := http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      n,
	}

	logger.Log.Info("Listening...")

	server.ListenAndServe()
}

// OnShutdown is called when the server has a panic.
func (app *AppServer) OnShutdown() {
	logger.OutputLog.Error("Executed OnShutdown")
}

// Special server handlers, outside of specific routes we have
func (app *AppServer) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Sender.JSON(w, http.StatusNotFound, fmt.Sprint("Not Found ", r.URL))
	if err != nil {
		panic(err)
	}
}

func (app *AppServer) NotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := app.Sender.JSON(w, http.StatusMethodNotAllowed, fmt.Sprint(r.Method, " method not allowed"))
	if err != nil {
		panic(err)
	}
}

// cSpell:ignore negroni httputils Nosniff urfave sirupsen logrus
