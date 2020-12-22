package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"

	envConfig "ec/config"
	"ec/config/cloudStorages"
	"ec/initializers"
	"ec/models"
	"ec/services/api/routes"
	"ec/utils"
)

func main() {
	initialize()

	e := echo.New()

	if envConfig.Env.Newrelic.AppName != "" && envConfig.Env.Newrelic.LicenseKey != "" {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName(envConfig.Env.Newrelic.AppName),
			newrelic.ConfigLicense(envConfig.Env.Newrelic.LicenseKey),
		)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		e.Use(nrecho.Middleware(app))
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(initializers.Auth)

	DefaultCORSConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	routes.SetWebInterfaces(e)
	routes.SetWsInterfaces(e)

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.HideBanner = true

	go func() {
		if err := e.Start(":3000"); err != nil {
			log.Println("start close echo")
			time.Sleep(500 * time.Millisecond)
			closeResource()
			log.Println("shutting down the server")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("accepted signal")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Println("shutting down failed, err:" + err.Error())
		e.Logger.Fatal(err)
	} else {
		log.Println("shutting down complete.")
	}
	closeResource()
}

func customHTTPErrorHandler(err error, context echo.Context) {
	language := context.Get("language").(string)
	if response, ok := err.(utils.Response); ok {
		response.Head["msg"] = fmt.Sprint(models.I18n.T(language, "error_code."+response.Head["code"]))
		context.JSON(http.StatusBadRequest, response)
	} else {
		panic(err)
	}
}

func initialize() {
	envConfig.InitEnv()
	cloudStorages.InitAwsS3Config()
	cloudStorages.InitQiniuConfig()

	initializers.InitMainDB()
	initializers.InitRedisPools()
	models.MainMigrations()
	models.InitI18n()
	// initializers.InitializeAmqpConfig()
	initializers.LoadInterfaces()
	// initializers.LoadCacheData()
	utils.SetLogAndPid()
}

func closeResource() {
	// initializers.CloseAmqpConnection()
	initializers.CloseRedisPools()
	initializers.CloseMainDB()
}
