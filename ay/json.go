package ay

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Json struct {
}

func (class Json) Msg(c *gin.Context, code string, msg string, data map[string]interface{}, args ...map[string]interface{}) {

	res := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	c.JSON(http.StatusOK, res)
	c.Abort()

}
