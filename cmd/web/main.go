package main

import ( 
    "os"
    "flag"
    "log/slog"
    "net/http"
    "database/sql"
    "html/template"
    "modulo.porreiro/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
    cfg           *config
    logger        *slog.Logger
    snippets      *models.SnippetModel
    templateCache map[string]*template.Template
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug, // Debug (descartadas), Info, Warn, Error
	}))

	var cfg config
	var dsn string

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network adress  ")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&dsn, "dsn", "web:pass@/snippetbox?parseTime=True", "MySWL data source name")
	flag.Parse() // Também existe flag.Int, flag.Bool...

	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

    defer db.Close()
    
    templateCache,err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    
    app := &application {
        cfg           : &cfg,
        logger        : logger,
        snippets      : &models.SnippetModel{ DB : db },
        templateCache : templateCache,
    }

	logger.Info("Starting server", slog.String("hosted_at", "https:://localhost"+app.cfg.addr))

	//func ListenAndServe(addr string, handler Handler) error
	err = http.ListenAndServe(app.cfg.addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, err
}
