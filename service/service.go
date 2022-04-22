package service

import (
	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine) *gin.Engine {
	r.Use(Cors())
	return r
}
