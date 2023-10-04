package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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
		log.Printf("Failed unmarshaling in addCats", err)
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

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	e.GET("/", yallo)

	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCat)

	e.POST("/dogs", addDog)

	e.POST("/hamsters", addHamster)

	e.Start(":8080")
}
