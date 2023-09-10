package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func newResponse(c *gin.Context, code int, message string, err error) {
	_, ok := err.(*gin.Error)
	if !ok {
		resp := response{
			Message: message,
			Error:   err.Error(),
		}
		c.AbortWithStatusJSON(code, resp)
		return
	}
	c.AbortWithStatusJSON(
		http.StatusInternalServerError, response{
			Message: "interval error",
			Error:   err.Error(),
		},
	)
	return
}
