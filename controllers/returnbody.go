package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReturnErrorBody(c *gin.Context, code int, msg string, err error) {
	data := make(map[string]string)
	data["err_log"] = fmt.Sprintf("%v", err)
	c.JSON(http.StatusBadRequest, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
