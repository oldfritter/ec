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
	newrelic "github.com/oldfritter/echo-middleware/v4"

	envConfig "ec/config"
	"ec/initializers"
	"ec/models"
	"ec/services/api/routes"
	"ec/utils"
)

func main() {
	initialize()

	e := echo.New()
	if envConfig.Env.Newrelic.AppName != "" && envConfig.Env.Newrelic.LicenseKey != "" {
		e.Use(newrelic.NewRelic(envConfig.Env.Newrelic.AppName, envConfig.Env.Newrelic.LicenseKey))
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(initializers.Auth)

	routes.SetWebInterfaces(e)
	routes.SetWsInterfaces(e)

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.HideBanner = true

	go func() {
		if err := e.Start(":9700"); err != nil {
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
