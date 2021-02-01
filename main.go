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
	httpSwagger "github.com/swaggo/http-swagger"
	"krungthai.com/khanapat/dpki/crypto-key-api/docs"
	_ "krungthai.com/khanapat/dpki/crypto-key-api/docs"
	"krungthai.com/khanapat/dpki/crypto-key-api/ecdsa"
	"krungthai.com/khanapat/dpki/crypto-key-api/key"
	"krungthai.com/khanapat/dpki/crypto-key-api/logger"
	"krungthai.com/khanapat/dpki/crypto-key-api/middleware"
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

// @title Crypto Key API
// @version 1.0
// @description API Service for managing key.
// @termsOfService http://swagger.io/terms/
// @contact.name KTB Blockchain Team
// @contact.url http://www.swagger.io/support
// @contact.email blockchain.info@krungthai.co.th
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /
// @schemes http https
func main() {
	route := mux.NewRouter()

	cfgCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                    // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE"}, // Allowing only get, just an example
		AllowedHeaders:   []string{"Content-Type", "Origin", "Authorization", "Accept"},
		AllowCredentials: true,
	})

	registerSwaggerRoute(route)

	logger := logger.NewLogConfig()

	middleware := middleware.NewMiddleware(logger)

	cryptoRoute := route.PathPrefix(viper.GetString("APP.CONTEXT.CRYPTO")).Subrouter()

	cryptoRoute.Use(middleware.JSONMiddleware)
	cryptoRoute.Use(middleware.ContextLogAndLoggingMiddleware)

	cryptoRoute.Handle("/ecdsa", ecdsa.NewAsymmetricEcdsaKey(
		ecdsa.NewGenerateEcdsaKeyFn(),
	)).Methods(http.MethodPost)

	cryptoRoute.Handle("/ecdsa/sign", ecdsa.NewSignEcdsaKey(
		ecdsa.NewSignEcdsaKeyFn(),
	)).Methods(http.MethodPost)

	cryptoRoute.Handle("/ecdsa/verify", ecdsa.NewVerifyEcdsaKey(
		ecdsa.NewVerifyEcdsaKeyFn(),
	)).Methods(http.MethodPost)

	cryptoRoute.Handle("/public_key/validate", key.NewValidationPublicKey(
		key.NewValidatePublicKeyFn(),
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

func registerSwaggerRoute(route *mux.Router) {
	route.PathPrefix("/crypto-key-api/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/crypto-key-api/swagger/doc.json", viper.GetString("APP.SWAGGER.HOST"))),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
	docs.SwaggerInfo.Host = viper.GetString("APP.SWAGGER.HOST")
	docs.SwaggerInfo.BasePath = viper.GetString("APP.CONTEXT.CRYPTO")
}
