package ginapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"hello": "world"})
}
