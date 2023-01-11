package common

import (
	"github.com/gin-gonic/gin"
)

func GetJsonBody[T any](c *gin.Context, body *T) error {
	err := c.BindJSON(&body)
	if err != nil {
		return err
	} else {
		return nil
	}
}
