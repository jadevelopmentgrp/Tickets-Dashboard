package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func MultiReadBody(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}
