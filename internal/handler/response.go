package handler

import (
	"TransactionsQueues/pkg"
	"errors"
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, err error) {
	var customError pkg.Error
	if errors.As(err, &customError) {
		c.AbortWithStatusJSON(customError.Code, customError.Description)
	} else {
		c.AbortWithStatusJSON(500, "Internal server error")
	}
}
