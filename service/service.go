package service

import (
	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine) *gin.Engine {
	r = SetSession(r)
	return r

}
