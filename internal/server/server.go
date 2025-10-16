package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ilhamgepe/gepay/pkg/config"
	"github.com/ilhamgepe/gepay/pkg/logger"
)

type Server struct {
	App    *gin.Engine
	Config config.App
	Log    logger.Logger
}

func NewServer(config config.App, log logger.Logger) *Server {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	redisHost := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	store, err := redis.NewStore(10, "tcp", redisHost, "", config.Redis.Password, []byte(config.Server.Secret))
	if err != nil {
		log.Fatal(err, "error connect redis")
		panic(err)
	}
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3600,
		Secure:   config.Server.Env == "production",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("gepay-session-hire-me", store))

	return &Server{
		App:    router,
		Config: config,
		Log:    log,
	}
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.Config.Server.Addr, s.Config.Server.Port),
		Handler: s.App,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Log.Fatalf(err, "failed to serve")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.Log.Infof("Server Shutdown: %v", err)
	}
	s.Log.Info("Server exiting")
}
