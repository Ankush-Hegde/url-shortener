package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"url-shortener/internal/pkg/server/gen"
	"url-shortener/internal/pkg/server/middleware"
	"url-shortener/internal/pkg/server/service"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux" // Use DataDog's package for metrics instead of the standard mux package
)

var (
	router *mux.Router
)

func setupRoutes() {
	urlShortenerImpl := service.NewUrlShortenAPIServiceImpl()
	urlshortenerController := gen.NewDefaultAPIController(urlShortenerImpl)

	routers := []gen.Router{
		urlshortenerController,
	}

	// use DD package for below line to enable metrics
	// router = mux.NewRouter(
	// 	mux.WithServiceName("url-shortener"),
	// 	mux.WithQueryParams(),
	// 	mux.WithAnalytics(true),
	// )

	router = mux.NewRouter()

	// register controllers to the router
	registerControllers(routers)

	// register Prometheus metrics endpoint
	router.Handle("/management/prometheus", promhttp.Handler())

	// -----swagger-ui-----
	openApiYamlDir := filepath.Join("api", "url-shortener")
	router.PathPrefix("/api/url-shortener/").Handler(http.StripPrefix("/api/url-shortener/", http.FileServer(http.Dir(openApiYamlDir))))

	swaggerDir := filepath.Join("third_party", "swagger-ui", "dist")
	router.HandleFunc("/doc/swagger-ui/swagger-initializer.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(openApiYamlDir, "swagger-initializer.js"))
	}).Methods(http.MethodGet)
	router.HandleFunc("/doc/swagger-ui", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/doc/swagger-ui/", http.StatusMovedPermanently)
	}).Methods(http.MethodGet)
	router.PathPrefix("/doc/swagger-ui/").Handler(
		http.StripPrefix("/doc/swagger-ui/", http.FileServer(http.Dir(swaggerDir))),
	)
	// ---------------------

	// middleware
	router.Use(middleware.ValidateRequestMiddleware)
}

func registerControllers(routers []gen.Router) {
	for _, r := range routers {
		for name, route := range r.Routes() {
			handler := route.HandlerFunc
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(name).
				Handler(handler)
		}
	}
}

func main() {
	// Configure slog to output structured JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// load env

	setupRoutes()
	port := 8089

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		slog.Error("Server failed to start", "err", err)
	}
}
