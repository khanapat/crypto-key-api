package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"krungthai.com/khanapat/dpki/generate-key-api/key"
	"krungthai.com/khanapat/dpki/generate-key-api/logger"
	"krungthai.com/khanapat/dpki/generate-key-api/middleware"
)

func init() {
	runtime.GOMAXPROCS(1)
	initViper()
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic("There is no such a config file in path ./config/config.yaml")
		} else {
			log.Panic("There is some problem about data in file")
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	route := mux.NewRouter()

	cfgCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                    // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE"}, // Allowing only get, just an example
		AllowedHeaders:   []string{"Content-Type", "Origin", "Authorization", "Accept"},
		AllowCredentials: true,
	})

	logger := logger.NewLogConfig()

	middleware := middleware.NewMiddleware(logger)

	keyRoute := route.PathPrefix(viper.GetString("APP.CONTEXT.KEY")).Subrouter()

	keyRoute.Use(middleware.AccessControl)
	keyRoute.Use(middleware.JsonHeader)
	keyRoute.Use(middleware.ContextLog)
	keyRoute.Use(middleware.LogRequestInfo)
	keyRoute.Use(middleware.LogResponseInfo)

	keyRoute.Handle("/ecdsa", key.NewAsymmetricEcdsaKey(
		key.NewGenerateEcdsaKeyFn(),
	)).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", viper.GetString("APP.PORT")),
		Handler:      cfgCors.Handler(route),
		ReadTimeout:  time.Duration(viper.GetInt("APP.TIMEOUT")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("APP.TIMEOUT")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("APP.TIMEOUT")) * time.Second,
	}

	logger.Info(fmt.Sprintf("⇨ http server started on [::]:%s", viper.GetString("APP.PORT")))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("APP.TIMEOUT"))*time.Second)
	defer cancel()

	srv.Shutdown(ctx)

	logger.Info("shutting down")
	os.Exit(0)
}
