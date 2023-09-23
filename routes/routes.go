package routes

import (
	"echoLearning/handlers"
	"echoLearning/structs"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, fakeDataBase []structs.Client) {
	e.POST("/client/:id/process/transaction", func(c echo.Context) error {
		return handlers.ProcessTransactionHandler(c, fakeDataBase)
	})
	e.POST("/client/:from/send/:to/:amount", func(c echo.Context) error {
		return handlers.SendFundsHandler(c, fakeDataBase)
	})
	e.POST("/client", func(c echo.Context) error {
		return handlers.CreateClientHandler(c, &fakeDataBase)
	})
	e.GET("/client/:id/balance", func(c echo.Context) error {
		return handlers.ClientBalanceHandler(c, fakeDataBase)
	})
	e.POST("/client/:id/transaction", func(c echo.Context) error {
		return handlers.CreateTransactionHandler(c, fakeDataBase)
	})
	e.GET("/clients", func(c echo.Context) error {
		return handlers.GetAllClientsHandler(c, fakeDataBase)
	})
}
