package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/money_transfer/internal/service"
	"net/http"
)

type accountRoutes struct {
	account service.Account
}

func newAccountRoutes(g *gin.RouterGroup, account service.Account) {
	r := &accountRoutes{
		account: account,
	}

	g.POST("/invoice", r.Deposit)
}

type requestBodyDeposit struct {
	CurrencyCode  string  `json:"currency_code"`
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"account_number"`
}

func (r *accountRoutes) Deposit(c *gin.Context) {
	var depositBody requestBodyDeposit

	if err := c.ShouldBindJSON(&depositBody); err != nil {
		newResponse(c, http.StatusBadRequest, "error parsing json body", err)
		return
	}

	input := service.AccountDepositInput{
		CurrencyCode: depositBody.CurrencyCode,
		Amount:       depositBody.Amount,
		To:           depositBody.AccountNumber,
	}

	err := r.account.Deposit(c, input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "error depositing", err)
		return
	}

	c.Status(http.StatusOK)
}
