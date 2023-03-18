package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
)

const matrixDatabase = "/tmp/matrix.db"

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.HideBanner = true

	eserver := new(Server)
	if err := envconfig.Process("app_server", eserver); err != nil {
		e.Logger.Fatal(err)
	}

	var authToken string
	authToken, err := readSecret(eserver.Http.AuthTokenFile)
	if err != nil {
		authToken = uuid.Must(uuid.NewV4()).String()
		fmt.Printf("Generated auth token is %s\n\n", authToken)
	}

	matrixHomeserverUrl, err := readSecret(eserver.Matrix.HomeserverUrlFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	matrixUserName, err := readSecret(eserver.Matrix.UserNameFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	matrixUserPassword, err := readSecret(eserver.Matrix.UserPasswordFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	matrixRoomID, err := readSecret(eserver.Matrix.RoomIdFile)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Matrix
	mcli, err := mautrix.NewClient(matrixHomeserverUrl, "", "")
	if err != nil {
		log.Fatal(err)
	}

	cryptoHelper, err := cryptohelper.NewCryptoHelper(mcli, []byte("meow"), matrixDatabase)
	if err != nil {
		log.Fatal(err)
	}

	cryptoHelper.LoginAs = &mautrix.ReqLogin{
		Type:       mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: matrixUserName},
		Password:   matrixUserPassword,
	}
	err = cryptoHelper.Init()
	if err != nil {
		log.Fatal(err)
	}
	mcli.Crypto = cryptoHelper

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 5 * time.Second,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: eserver.Http.CorsAllowOrigins,
		AllowMethods: []string{http.MethodPost},
	}))

	// Handlers
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.POST("/", handleEvent(mcli, matrixRoomID), middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == authToken, nil
	}))

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", eserver.Http.BindPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := cryptoHelper.Close(); err != nil {
		e.Logger.Fatal("Error closing database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func readSecret(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(string(b))

	return result, nil
}
