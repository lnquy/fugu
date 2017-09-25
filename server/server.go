package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lnquy/fugu/config"
	"github.com/lnquy/fugu/modules/util"
	"github.com/lnquy/fugu/router"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func Serve(cfg *config.Config) {
	r := chi.NewRouter()
	configureRouter(cfg, r)

	addr := cfg.Server.GetFullAddr()
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.RTimeout,
		WriteTimeout: cfg.Server.WTimeout,
	}
	if cfg.Runtime.IsDebugging { // Only allow CORS in development mode
		srv.Handler = cors.Default().Handler(r)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill)

	go func() {
		if cfg.Server.TLSCert != "" && cfg.Server.TLSKey != "" {
			log.Infof("main: serving at https://%s", addr)
			log.Fatal(srv.ListenAndServeTLS(cfg.Server.TLSCert, cfg.Server.TLSKey))
		} else {
			log.Infof("main: serving at http://%s", addr)
			log.Fatal(srv.ListenAndServe())
		}
	}()

	// Graceful shutdown
	<-stop
	log.Info("main: shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	srv.Shutdown(ctx)
	log.Info("main: have a nice day, goodbye!")
}

func configureRouter(cfg *config.Config, r *chi.Mux) {
	if cfg.Runtime.IsDebugging {
		r.Use(middleware.DefaultLogger)
	}
	//r.Use(mdw.BrowserCache) // TODO: Cache for static files only
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Recoverer)

	// Routing
	dir := path.Join(util.GetWd(), "templates", "static")
	fileServer(r, "/static", http.Dir(dir))
	r.Get("/", router.GetIndexPage)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/fugu", func(r chi.Router) {
			r.Post("/lang/{lang}/arch/{arch}", router.CalcSizeOfStruct)
			r.Post("/lang/{lang}/arch/{arch}/optimize", router.OptimizeStruct)
		})
	})
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	//if strings.ContainsAny(path, ":*") {
	//	panic("FileServer does not permit URL parameters.")
	//}
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
