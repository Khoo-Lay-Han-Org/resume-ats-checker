package api

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RunOpenAPIDoc() *echo.Echo {
	router := echo.New()

	router.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	return router
}
