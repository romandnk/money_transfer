package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/money_transfer/internal/models"
	"github.com/romandnk/money_transfer/internal/service"
	"net/http"
)

type userRoutes struct {
	user service.User
}

func newUserRoutes(g *gin.RouterGroup, user service.User) {
	r := &userRoutes{
		user: user,
	}

	g.POST("/sign-in", r.SignIn)
	g.POST("/sign-up", r.SignUp)
}

type requestUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *userRoutes) SignUp(c *gin.Context) {
	var userBody requestUserBody

	if err := c.ShouldBindJSON(&userBody); err != nil {
		newResponse(c, http.StatusBadRequest, "error parsing json body", err)
		return
	}

	user := models.User{
		Email:    userBody.Email,
		Password: userBody.Password,
	}

	id, err := r.user.SignUp(c, user)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "error signing up", err)
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func (r *userRoutes) SignIn(c *gin.Context) {
	var userBody requestUserBody

	if err := c.ShouldBindJSON(&userBody); err != nil {
		newResponse(c, http.StatusBadRequest, "error parsing json body", err)
		return
	}

	token, err := r.user.SignIn(c, userBody.Email, userBody.Password)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "error signing up", err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": token})
}
