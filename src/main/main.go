package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from the web side!")
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	e.GET("/", Handler)

	e.Start(":8080")
}
