package testing

import (
	"echoLearning/handlers"
	"echoLearning/structs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllClientsHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/clients")

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	if assert.NoError(t, handlers.GetAllClientsHandler(c, fakeDataBase)) {

		assert.Equal(t, http.StatusOK, rec.Code)

		expectedJSON := `[{"clientID":1,"balance":10000},{"clientID":2,"balance":10000}]`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}

func TestPositiveClientBalanceHandler(t *testing.T) {

	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	var fakeDataBase []structs.Client = []structs.Client{*client1}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:id/balance")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handlers.ClientBalanceHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedJSON := `{"balance":10000, "clientID":1}`
		assert.JSONEq(t, expectedJSON, rec.Body.String())
	}
}

func TestNegativeClientBalanceHandler(t *testing.T) {

	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	var fakeDataBase []structs.Client = []structs.Client{*client1}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:id/balance")
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, handlers.ClientBalanceHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		expectedString := "\"client is not in massive\"\n"
		assert.Equal(t, expectedString, rec.Body.String())
	}
}

func TestPositiveProcessTransactionHandler(t *testing.T) {

	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	var fakeDataBase []structs.Client = []structs.Client{*client1}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:id/process/transaction")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handlers.ProcessTransactionHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestNegativeProcessTransactionHandler(t *testing.T) {

	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	var fakeDataBase []structs.Client = []structs.Client{*client1}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:id/process/transaction")
	c.SetParamNames("id")
	c.SetParamValues("2")

	if assert.NoError(t, handlers.ProcessTransactionHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestPositiveCreateClientHandler(t *testing.T) {
	e := echo.New()

	var fakeDataBase []structs.Client

	req := httptest.NewRequest(http.MethodPost, "/client?id=1&balance=1000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateClientHandler(c, &fakeDataBase)) {

		assert.Equal(t, http.StatusCreated, rec.Code)

		assert.Len(t, fakeDataBase, 1)

	}
}

func TestNegative1CreateClientHandler(t *testing.T) {
	e := echo.New()

	var fakeDataBase []structs.Client

	req := httptest.NewRequest(http.MethodPost, "/client?id=1&balance=-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateClientHandler(c, &fakeDataBase)) {

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expectedString := "\"incorrect balance\"\n"
		assert.Equal(t, expectedString, rec.Body.String())
	}
}

func TestNegative2CreateClientHandler(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	var fakeDataBase []structs.Client = []structs.Client{*client1}

	req := httptest.NewRequest(http.MethodPost, "/client?id=1&balance=1000", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateClientHandler(c, &fakeDataBase)) {

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expectedString := "\"Client ID already exists\"\n"
		assert.Equal(t, expectedString, rec.Body.String())

	}
}

func TestPositiveSendFundsHandler(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 5000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:from/send/:to/:amount")
	c.SetParamNames("from", "to", "amount")
	c.SetParamValues("1", "2", "100")

	if assert.NoError(t, handlers.SendFundsHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		assert.Equal(t, 9900, client1.Account.Balance)
		assert.Equal(t, 5100, client2.Account.Balance)
	}
}

func TestNegative1SendFundsHandler(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 5000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:from/send/:to/:amount")
	c.SetParamNames("from", "to", "amount")
	c.SetParamValues("1", "1", "100")

	if assert.NoError(t, handlers.SendFundsHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusConflict, rec.Code)

		expectedString := "\"sender ID equal to receiver ID\"\n"
		assert.Equal(t, expectedString, rec.Body.String())
	}
}

func TestNegative2SendFundsHandler(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 5000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:from/send/:to/:amount")
	c.SetParamNames("from", "to", "amount")
	c.SetParamValues("1", "3", "100")

	if assert.NoError(t, handlers.SendFundsHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expectedString := "\"receiver ID not found\"\n"
		assert.Equal(t, expectedString, rec.Body.String())
	}
}

func TestNegative3SendFundsHandler(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 5000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:from/send/:to/:amount")
	c.SetParamNames("from", "to", "amount")
	c.SetParamValues("3", "1", "100")

	if assert.NoError(t, handlers.SendFundsHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expectedString := "\"sender ID not found\"\n"
		assert.Equal(t, expectedString, rec.Body.String())
	}
}

func TestNegative4SendFundsHandler(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	bankAccountC2 := &structs.BankAccount{Balance: 5000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction)}
	client2 := &structs.Client{ID: 2, Account: bankAccountC2, TransChan: make(chan structs.Transaction)}

	var fakeDataBase []structs.Client = []structs.Client{*client1, *client2}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/client/:from/send/:to/:amount")
	c.SetParamNames("from", "to", "amount")
	c.SetParamValues("3", "1", "-100")

	if assert.NoError(t, handlers.SendFundsHandler(c, fakeDataBase)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expectedString := "\"incorrect balance\"\n"
		assert.Equal(t, expectedString, rec.Body.String())
	}
}

func TestPositiveCreateTransactionHandlerWithProcessTransaction(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction, 1)}
	client1.TransBool = true

	fakeDataBase := []structs.Client{*client1}

	reqProcess := httptest.NewRequest(http.MethodPost, "/client/:id/process/transaction", nil)
	recProcess := httptest.NewRecorder()
	cProcess := e.NewContext(reqProcess, recProcess)
	cProcess.SetParamNames("id")
	cProcess.SetParamValues("1")

	err := handlers.ProcessTransactionHandler(cProcess, fakeDataBase)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	assert.Equal(t, http.StatusOK, recProcess.Code)

	reqCreate := httptest.NewRequest(http.MethodPost, "/client/:id/transaction?amount=100&isDebit=true", nil)
	recCreate := httptest.NewRecorder()
	cCreate := e.NewContext(reqCreate, recCreate)
	cCreate.SetParamNames("id")
	cCreate.SetParamValues("1")

	if assert.NoError(t, handlers.CreateTransactionHandler(cCreate, fakeDataBase)) {
		assert.Equal(t, http.StatusOK, recCreate.Code)

		assert.Equal(t, 9900, client1.Account.Balance)
	}
}

func TestNegativeCreateTransactionHandlerWithProcessTransaction(t *testing.T) {
	e := echo.New()

	bankAccountC1 := &structs.BankAccount{Balance: 10000}
	client1 := &structs.Client{ID: 1, Account: bankAccountC1, TransChan: make(chan structs.Transaction, 1)}

	fakeDataBase := []structs.Client{*client1}

	reqCreate := httptest.NewRequest(http.MethodPost, "/client/:id/transaction?amount=100&isDebit=true", nil)
	recCreate := httptest.NewRecorder()
	cCreate := e.NewContext(reqCreate, recCreate)
	cCreate.SetParamNames("id")
	cCreate.SetParamValues("1")

	if assert.NoError(t, handlers.CreateTransactionHandler(cCreate, fakeDataBase)) {
		assert.Equal(t, http.StatusConflict, recCreate.Code)
		expectedString := "\"ProcessTransaction is not activated\"\n"
		assert.Equal(t, expectedString, recCreate.Body.String())
	}
}
