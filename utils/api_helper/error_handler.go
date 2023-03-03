package api_helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandlerError 错误处理
func HandlerError(g *gin.Context, err error) {
	g.JSON(
		http.StatusBadRequest, ErrResponse{
			Message: err.Error(),
		})
	g.Abort()
	return
}
