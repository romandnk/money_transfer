package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		h.logger.Info("Request info HTTP",
			zap.String("client ip", c.RemoteIP()),
			zap.String("method", c.Request.Method),
			zap.String("method path", c.FullPath()),
			zap.String("HTTP version", c.Request.Proto),
			zap.Int("status code", c.Writer.Status()),
			zap.String("processing time", duration.String()),
		)
	}
}

func (h *Handler) authorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				return nil, errors.New("error parsing token")
			}
			return []byte(h.signKey), nil
		})
		if err != nil || !token.Valid {
			newResponse(c, http.StatusUnauthorized, "unauthorized", err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["user_id"].(string)
			c.Set("user_id", userId)
			c.Next()
		}
	}
}
