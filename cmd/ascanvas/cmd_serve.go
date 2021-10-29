package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/broadcaster/memory"
	"github.com/fluxynet/ascanvas/cmd"
	docs "github.com/fluxynet/ascanvas/docs/ascanvas"
	"github.com/fluxynet/ascanvas/internal"
	"github.com/fluxynet/ascanvas/repo/sequel"
	"github.com/fluxynet/ascanvas/web"
	"github.com/fluxynet/ascanvas/web/canvas"
)

func Serve(_ *cobra.Command, _ []string) {
	var (
		logger        *zap.Logger
		db            *sql.DB
		repo          *sequel.Repository
		broadcaster   *memory.Memory
		canvasService *ascanvas.CanvasService
		webCanvas     canvas.WebCanvas
		router        *chi.Mux

		config = Config{
			ListenAddr: "127.0.0.1:1337",
			DbDriver:   "sqlite",
			DSN:        "ascanvas.db",
			LogLevel:   "debug",
		}

		err = cmd.LoadConfig("ascanvas.json", &config)
	)

	if err != nil {
		log.Printf("config file not loaded, using defaults (%s)\n", err.Error())
	}

	logger, err = cmd.Logger(config.LogLevel, cmd.DoNotLogToFile)

	if err != nil {
		log.Fatalln("failed to start logger: ", err.Error())
	}

	broadcaster = memory.New()
	defer internal.Closed(broadcaster)

	db, err = sql.Open(config.DbDriver, config.DSN)
	if err != nil {
		log.Fatalln("failed to open database connection: ", err.Error())
	} else if err = db.Ping(); err != nil {
		log.Fatalln("failed to ping database: ", err.Error())
	} else if _, err = db.Exec(sequel.SQLiteSchemaInit); err != nil {
		log.Fatalln("failed to initialize schema: ", err.Error())
	}

	repo = &sequel.Repository{DB: db}

	canvasService = &ascanvas.CanvasService{
		Repo:        repo,
		BroadCaster: broadcaster,
		Logger:      logger,
		GenerateID:  ascanvas.UUIDGenerator,
		Broadcast:   ascanvas.AsyncBroadcast,
	}

	webCanvas = canvas.WebCanvas{
		Service: canvasService,
		GetID:   web.ChiIDGetter,
	}

	docs.SwaggerInfo.Host = config.ListenAddr
	docs.SwaggerInfo.BasePath = "/api"

	router = chi.NewMux()
	router.Get("/swagger/*", httpSwagger.Handler())

	router.Route("/api", func(r chi.Router) {
		r.Get("/events", webCanvas.Observe)

		r.Get("/{id}/events", webCanvas.Observe)
		r.Patch("/{id}/rectangle", webCanvas.Rectangle)
		r.Patch("/{id}/floodfill", webCanvas.Floodfill)
		r.Delete("/{id}", webCanvas.Delete)
		r.Get("/{id}", webCanvas.Get)

		r.Post("/", webCanvas.Create)
		r.Get("/", webCanvas.List)
	})

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", web.ContentTypeHTML)
		w.WriteHeader(http.StatusOK)
		w.Write(ascanvas.IndexHTML)
	})

	log.Printf("server listening at %s\n", config.ListenAddr)
	err = http.ListenAndServe(config.ListenAddr, router)
	if err != nil {
		log.Fatalln(err)
	}
}
