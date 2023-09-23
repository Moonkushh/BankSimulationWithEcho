package main

import (
	"echoLearning/routes"
	"echoLearning/structs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func main() {

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	routes.RegisterRoutes(e, fakeDataBase)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Start(":1488")
}
