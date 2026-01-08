package main 

import ( 
    "os"
    "flag"
    "log/slog"
    "net/http"
)


type config struct {
    addr      string
    staticDir string
}

type application struct {
    logger *slog.Logger
    cfg    *config
}


func main () {

    logger := slog.New ( slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,
        Level : slog.LevelDebug,
    // Debug (descartadas), Info, Warn, Error 
    }))
    
    var cfg config
    
    flag.StringVar( &cfg.addr      ,"addr"      , ":4000"       , "HTTP network adress  ") 
    flag.StringVar( &cfg.staticDir ,"static-dir", "./ui/static/", "Path to static assets") 
    flag.Parse()
    // Também existe flag.Int, flag.Bool...

    app := &application {
        logger: logger,
        cfg   : &cfg,
    }

    logger.Info("Starting server" , slog.String( "hosted_at", "https:://localhost" + app.cfg.addr))

    //func ListenAndServe(addr string, handler Handler) error
    err := http.ListenAndServe( app.cfg.addr, app.routes())

    logger.Error(err.Error())
    os.Exit(1)
}
