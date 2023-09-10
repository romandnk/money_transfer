package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/money_transfer/internal/service"
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

func (r *accountRoutes) Deposit(c *gin.Context) {

}
