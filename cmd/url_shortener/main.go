package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"url-shortener/config"
	"url-shortener/internal/pkg/database/mongodb"
	"url-shortener/internal/pkg/redis"
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

func setupRoutes(mongoClient *mongodb.Client, redisClient *redis.Client) {
	urlShortenerImpl := service.NewUrlShortenAPIServiceImpl(mongoClient, redisClient)
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
	router.Use(middleware.CorsMiddleware)
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

func initMongoDB(ctx context.Context, mongoURI string) (*mongodb.Client, error) {
	mongoClient, err := mongodb.Connect(ctx, mongoURI)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}
	slog.Info("Connected to MongoDB")
	return mongoClient, nil
}

func initRedis(ctx context.Context, addr, username, password string, db int, tls bool) (*redis.Client, error) {
	redisClient, err := redis.NewClient(ctx, redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
		TLS:      tls,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to Redis: %w", err)
	}
	slog.Info("Connected to Redis")
	return redisClient, nil
}

func runServer() error {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("load environment config: %w", err)
	}

	mongoClient, err := initMongoDB(context.Background(), cfg.MongoURI)
	if err != nil {
		return err
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			slog.Error("Failed to disconnect from MongoDB", "error", err)
		}
	}()

	redisClient, err := initRedis(context.Background(), cfg.RedisAddr, cfg.RedisUsername, cfg.RedisPassword, cfg.RedisDB, cfg.RedisTLS)
	if err != nil {
		return err
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("Failed to close Redis client", "error", err)
		}
	}()

	setupRoutes(mongoClient, redisClient)

	address := fmt.Sprintf(":%s", cfg.ServerPort)
	slog.Info("Starting URL shortener server", "address", address)

	if err := http.ListenAndServe(address, router); err != nil {
		return fmt.Errorf("server failed on %s: %w", address, err)
	}
	return nil
}

func main() {
	// Configure slog to output structured JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := runServer(); err != nil {
		slog.Error("URL shortener exited", "error", err)
		os.Exit(1)
	}
}
