package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Root(c *gin.Context) {
	c.HTML(http.StatusOK, "root/index.tmpl", gin.H{
		"title": "Posts2",
	})
}
