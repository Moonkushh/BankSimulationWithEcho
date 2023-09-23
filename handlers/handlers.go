package handlers

import (
	"echoLearning/operations"
	"echoLearning/structs"
	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"strconv"
)

// ProcessTransactionHandler @Summary Processing a customer transaction
// @Description Processing a customer transaction by its ID
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /client/{id}/process/transaction [post]
func ProcessTransactionHandler(c echo.Context, fakeDataBase []structs.Client) error {
	clientIDStr := c.Param("id")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid client ID")
	}

	var client *structs.Client

	for i := range fakeDataBase {
		if fakeDataBase[i].ID == clientID {
			client = &fakeDataBase[i]
		}
	}

	if client == nil {
		return c.JSON(http.StatusNotFound, "Client not found")
	}

	client.TransBool = true

	go func() {
		err := operations.ProcessTransactions(*client)
		if err != nil {
			log.Printf("error %s", err)
		}
		close(client.TransChan)
	}()
	return c.NoContent(http.StatusOK)
}

// SendFundsHandler @Summary Sending funds between clients
// @Description Sending funds from one client to another
// @Accept json
// @Produce json
// @Param from path int true "Sender ID"
// @Param to path int true "Receiver ID"
// @Param amount path int true "Transaction object"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /client/{from}/send/{to}/{amount} [post]
func SendFundsHandler(c echo.Context, fakeDataBase []structs.Client) error {
	amountStr := c.Param("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid amount")
	}

	if amount < 0 {
		return c.JSON(http.StatusBadRequest, "incorrect balance")
	}

	fromIDStr := c.Param("from")
	toIDStr := c.Param("to")

	fromID, err := strconv.Atoi(fromIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid sender's ID")
	}

	toID, err := strconv.Atoi(toIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid recipient's ID")
	}

	var sender, receiver *structs.Client
	for i := range fakeDataBase {
		if fakeDataBase[i].ID == fromID {
			sender = &fakeDataBase[i]
		}
		if fakeDataBase[i].ID == toID {
			receiver = &fakeDataBase[i]
		}
	}

	if sender == receiver {
		return c.JSON(http.StatusConflict, "sender ID equal to receiver ID")
	}

	if sender == nil {
		return c.JSON(http.StatusBadRequest, "sender ID not found")
	} else if receiver == nil {
		return c.JSON(http.StatusBadRequest, "receiver ID not found")
	}

	err = operations.SendFunds(sender, receiver, amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// ClientBalanceHandler @Summary Getting the client's balance
// @Description Getting current client balance by client ID
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} interface{} "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /client/{id}/balance [get]
func ClientBalanceHandler(c echo.Context, fakeDataBase []structs.Client) error {
	clientIDStr := c.Param("id")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid client ID")
	}

	var client *structs.Client

	for i := range fakeDataBase {
		if clientID == fakeDataBase[i].ID {
			client = &fakeDataBase[i]
		}
	}

	if client == nil {
		return c.JSON(http.StatusNotFound, "client is not in massive")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"clientID": client.ID,
		"balance":  client.Account.Balance,
	})
}

// CreateClientHandler @Summary Creating a new client
// @Description Create a new client with the specified ID and add it to the database
// @Accept json
// @Produce json
// @Param id query int true "Client ID"
// @Param balance query int true "Client balance"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Router /client [post]
func CreateClientHandler(c echo.Context, fakeDataBase *[]structs.Client) error {
	idStr := c.QueryParam("id")
	balanceStr := c.QueryParam("balance")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid client ID")
	}

	balance, err := strconv.Atoi(balanceStr)

	if balance < 0 {
		return c.JSON(http.StatusBadRequest, "incorrect balance")
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid balance")
	}

	for _, client := range *fakeDataBase {
		if client.ID == id {
			return c.JSON(http.StatusBadRequest, "Client ID already exists")
		}
	}

	newClient := &structs.Client{
		ID:        id,
		Account:   &structs.BankAccount{Balance: balance},
		TransChan: make(chan structs.Transaction),
	}

	*fakeDataBase = append(*fakeDataBase, *newClient)

	return c.String(http.StatusCreated, "Created")
}

// CreateTransactionHandler @Summary Creating a new transaction
// @Description Create a new transaction and add it to the client's transaction channel
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param amount query int true "Transaction amount"
// @Param isDebit query bool true "Is Debit transaction (true for debit, false for credit)"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "StatusConflict"
// @Router /client/{id}/transaction [post]
func CreateTransactionHandler(c echo.Context, fakeDataBase []structs.Client) error {
	clientIDStr := c.Param("id")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid client ID")
	}

	amountStr := c.QueryParam("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid amount")
	}

	isDebitStr := c.QueryParam("isDebit")
	isDebit, err := strconv.ParseBool(isDebitStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid isDebit value")
	}

	var client *structs.Client

	for i := range fakeDataBase {
		if fakeDataBase[i].ID == clientID {
			client = &fakeDataBase[i]
		}
	}

	if client == nil {
		return c.JSON(http.StatusNotFound, "Client not found")
	}

	if client.TransBool != true {
		return c.JSON(http.StatusConflict, "ProcessTransaction is not activated")
	}

	transaction := structs.Transaction{
		Amount:  amount,
		IsDebit: isDebit,
	}

	log.Printf("Sending transaction %+v to client %d", transaction, client.ID)

	go operations.PerformTransactions(client, transaction)

	return c.NoContent(http.StatusOK)
}

// GetAllClientsHandler @Summary Getting information about all clients
// @Description Get information about all clients including their IDs and balances
// @Accept json
// @Produce json
// @Success 200 {object} interface{} "OK"
// @Router /clients [get]
func GetAllClientsHandler(c echo.Context, fakeDataBase []structs.Client) error {
	clientsInfo := []map[string]interface{}{}

	for _, client := range fakeDataBase {
		clientInfo := map[string]interface{}{
			"clientID": client.ID,
			"balance":  client.Account.Balance,
		}
		clientsInfo = append(clientsInfo, clientInfo)
	}

	return c.JSON(http.StatusOK, clientsInfo)
}
