package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from the web side!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Your cat name is: %s\nAnd his type is %s\n", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to lets us know if you want json or string data",
	})

}

// most fast
func addCat(c echo.Context) error {
	cat := Cat{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)

	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("This is your cat: %v", cat)
	return c.String(http.StatusOK, "We got your cat!")
}

func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)

	if err != nil {
		log.Printf("Failed proccesing addDog request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your dog: %v", dog)
	return c.String(http.StatusOK, "We got your dog!")

}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)

	if err != nil {
		log.Printf("Failed proccesing addHamster request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your hamster: %v", hamster)
	return c.String(http.StatusOK, "We got your hamster!")
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "Horay you are on the secret admin main page!")
}

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "You are on the secret cookie page")
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// check username and password against DB after hashing the password
	if username == "jack" && password == "1234" {
		cookie := &http.Cookie{}

		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)

		c.SetCookie(cookie)

		return c.String(http.StatusOK, "You were logged on!")

	}

	return c.String(http.StatusOK, "Your username or password were wrong")
}

// /////////////////////////// middlewares /////////////////////////////
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "rovezuka/1.0")
		c.Response().Header().Set("MyServer", "This is my server")

		return next(c)
	}
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")

		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.String(http.StatusUnauthorized, "you dont have any cookie")
			}
			log.Println(err)
			return err
		}

		if cookie.Value == "some_string" {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "You dont have the right cookie")
	}
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	e.Use(ServerHeader)

	adminGroup := e.Group("/admin")

	cookieGroup := e.Group("/cookie")

	// this logs the server interaction
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}  ${status}  ${method}  ${host}${path}  ${latency_human}]" + "\n",
	}))

	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// check in the DB
		if username == "jack" && password == "1234" {
			return true, nil
		}

		return false, nil

	}))

	cookieGroup.Use(checkCookie)

	cookieGroup.GET("/main", mainCookie)

	adminGroup.GET("/main", mainAdmin)

	e.GET("login", login)

	e.GET("/", yallo)

	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCat)

	e.POST("/dogs", addDog)

	e.POST("/hamsters", addHamster)

	e.Start(":8080")
}
